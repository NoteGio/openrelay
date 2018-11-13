/*
This module serves the security, rather than the function, of OpenRelay.

The 0x Protocol uses a TokenTransferProxy contract to transfer tokens between
users. Users set an allowance for the TokenTransferProxy, giving it permission
to swap tokens on their behalf. Users can verify that the TokenTransferProxy
will only allow authorized exchanges to transfer tokens under pre-defined
conditions.

The TokenTransferProxy is managed by a timelocked multisig contract, allowing
the 0x team to add new exchange contracts to the TokenTransferProxy. If they
want to add a new exchange, they must submit a transaction to the multisig
contract, submit the necessary number of confirmations, then wait 14 days for
the timelock to allow the transaction to be submitted to the
TokenTransferProxy.

If the 0x team ever decided to do anything malicious, users would have 14 days
to retract their allowances before the malicious contract was authorized by the
TokenTransferProxy. This service watches the multisig contract and logs such
events. The logs in turn become alerts to help ensure that any malicious
submissions were detected immediately, and give users as much time as possible
to withdraw their allowances.

*/

package multisig

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/monitor/blocks"
	"log"
)

type multisigBlockConsumer struct {
	multisigAddress   *big.Int
	logFilter         ethereum.LogFilterer
}

func (consumer *multisigBlockConsumer) Consume(delivery channels.Delivery) {
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	if types.BloomLookup(block.Bloom, consumer.multisigAddress) {
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: []common.Address{common.BigToAddress(consumer.multisigAddress)},
			Topics: [][]common.Hash{
				nil,
				nil,
				nil,
			},
		}
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting: %v", block.Number, err.Error())
		}
		if len(logs) > 0 {
			log.Printf("Multsig Contract '%v' Received Event in block %#x", common.BigToAddress(consumer.multisigAddress), block.Hash[:])
		}
	}
	delivery.Ack()
}

func NewMultisigBlockConsumer(multisigAddress *big.Int, lf ethereum.LogFilterer) (channels.Consumer) {
	submissionTopic := &big.Int{}
	submissionTopic.SetString("c0ba8fe4b176c1714197d43b9cc6bcf797a4a7461c5fe8d0ef6e184ae7601e51", 16)
	return &multisigBlockConsumer{multisigAddress, lf}
}

func NewRPCMultisigBlockConsumer(rpcURL string, multisigAddress string) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return NewMultisigBlockConsumer(common.HexToAddress(multisigAddress).Big(), client), nil
}
