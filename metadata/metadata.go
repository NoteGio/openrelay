package metadata

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	orCommon "github.com/notegio/openrelay/common"
	tokenModule "github.com/notegio/openrelay/token"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	dbModule "github.com/notegio/openrelay/db"
	"strings"
	"io/ioutil"
)

type TokenURIGetter interface {
	TokenURI(_tokenId *big.Int) (string, error)
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type OrderMetadataConsumer struct {
	conn   bind.ContractBackend
	client HttpClient
	db     *gorm.DB
	s      orCommon.Semaphore
}

func (consumer *OrderMetadataConsumer) ProcessAssetData(data *types.AssetData) {
	if data.IsType(types.ERC721ProxyID) {
		token, err := tokenModule.NewERC721Token(orCommon.ToGethAddress(data.Address()), consumer.conn)
		uri, err := token.TokenURI(nil, data.TokenID().Big())
		if err != nil {
			log.Printf("Error getting token URI for asset %#x: '%v'", (*data)[:], err.Error())
			return
		}
		if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
			resp, err := consumer.client.Get(uri)
			if err != nil {
				log.Printf("Error resolving token URI for asset %#x (%v): '%v' - %v", (*data)[:], uri, err.Error(), resp)
				// TODO: Include URI in metadata even if it doesn't resolve when we try it
				return
			}
			jsonData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error getting response data for asset %#x (%v): '%v'", (*data)[:], uri, err.Error())
				return
			}
			metadata := &dbModule.AssetMetadata{}
			if err := json.Unmarshal(jsonData, &metadata); err != nil {
				log.Printf("Error parsing JSON for asset %#x (%v): '%v'", (*data)[:], string(jsonData), err.Error())
				metadata.RawMetadata = string(jsonData)
				if len(metadata.RawMetadata) > 512 {
					metadata.RawMetadata = metadata.RawMetadata[:512]
				}
				// Don't return, because we still want to save the asset with the raw
				// metadata in case that's useful to anyone.
			}
			metadata.URI = uri
			metadata.SetAssetData(*data)
			if err := consumer.db.Model(&dbModule.AssetMetadata{}).Save(metadata).Error; err != nil {
				log.Printf("Error saving metdata for asset %#x: %v", (*data)[:], err.Error())
			}

		} else {
			// Eventually we may add support for ipfs://, swarm://, and others
			log.Printf("Unknown URI scheme: %v", uri)
		}
	}
}

func (consumer *OrderMetadataConsumer) Consume(msg channels.Delivery) {
	consumer.s.Acquire()
	go func(){
		defer consumer.s.Release()
		order, err := types.OrderFromBytes([]byte(msg.Payload()))
		if err != nil {
			log.Printf("Failed to parse order: %v", err.Error())
			msg.Reject()
			return
		}
		consumer.ProcessAssetData(&order.MakerAssetData)
		consumer.ProcessAssetData(&order.TakerAssetData)
		msg.Ack()
	}()
}
func NewOrderMetadataConsumer(rpcURL string, db *gorm.DB, concurrency int) (*OrderMetadataConsumer, error){
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	if _, err = conn.SyncProgress(context.Background()); err != nil {
		// This is just here so that an RpcBalanceChecker can't be instantiated
		// successfully if the RPC server isn't responding properly. What RPC
		// function we call isn't important, but SyncProgress is pretty cheap.
		return nil, err
	}
	consumer := &OrderMetadataConsumer{
		conn: conn,
		client: http.DefaultClient,
		db: db,
		s: orCommon.NewSemaphore(concurrency),
	}
	return consumer, nil
}
func NewRawOrderMetadataConsumer(conn bind.ContractBackend, client HttpClient, db *gorm.DB, concurrency int) (*OrderMetadataConsumer, error) {
	return &OrderMetadataConsumer{conn, client, db, orCommon.NewSemaphore(concurrency)}, nil
}
