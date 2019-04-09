package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"math/big"
	"time"
	"errors"
)

const (
	StatusOpen      = int64(0)
	StatusFilled    = int64(1)
	StatusUnfunded  = int64(2)
	StatusCancelled = int64(3)
)

func DefaultSha3() []byte {
	// sha3("nil")
	return []byte{167, 216, 239, 244, 2, 111, 37, 45, 181, 185, 12, 120, 228, 61, 209, 145, 223, 230, 229, 95, 203, 152, 84, 138, 95, 56, 250, 240, 212, 227, 235, 57}
}

type Order struct {
	types.Order
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderHash []byte  `gorm:"primary_key"`
	Status    int64   `gorm:"index"`
	Price     float64 `gorm:"index:price"`
	FeeRate   float64 `gorm:"index:price"`
	MakerAssetRemaining *types.Uint256
	MakerFeeRemaining   *types.Uint256
	MakerAssetMetadata  *AssetMetadata
	TakerAssetMetadata  *AssetMetadata
}

func (order *Order) TableName() string {
	return "orderv2"
}

func (order *Order) Populate() {
	order.OrderHash = order.Hash()
	remainingAmount := order.TakerAssetAmount.Big()
	thresholdAmount := order.TakerAssetAmount.Big()
	// If the order is larger than 100 base units, we want to remove it from the
	// orderbook when it's 99% filled. For fewer than 100 base units, integer
	// arithmetic makes this not work the way we want it to.
	if thresholdAmount.Cmp(big.NewInt(100)) > 0 {
		thresholdAmount.Mul(thresholdAmount, big.NewInt(99))
		thresholdAmount.Div(thresholdAmount, big.NewInt(100))
	}
	remainingAmount.Sub(remainingAmount, order.TakerAssetAmountFilled.Big())
	thresholdAmount.Sub(thresholdAmount, order.TakerAssetAmountFilled.Big())

	// makerRemainingInt = (MakerAssetAmount * RemainingAmount) / TakerAssetAmount
	makerRemainingInt := new(big.Int).Div(new(big.Int).Mul(order.MakerAssetAmount.Big(), remainingAmount), order.TakerAssetAmount.Big())
	makerRemaining := &types.Uint256{}
	copy(makerRemaining[:], abi.U256(makerRemainingInt))
	makerFeeRemainingInt := new(big.Int).Div(new(big.Int).Mul(order.MakerFee.Big(), remainingAmount), order.TakerAssetAmount.Big())
	makerFeeRemaining := &types.Uint256{}
	copy(makerFeeRemaining[:], abi.U256(makerFeeRemainingInt))
	order.MakerAssetRemaining = makerRemaining
	order.MakerFeeRemaining = makerFeeRemaining

	takerTokenAmount := new(big.Float).SetInt(order.TakerAssetAmount.Big())
	makerTokenAmount := new(big.Float).SetInt(order.MakerAssetAmount.Big())
	takerFeeAmount := new(big.Float).SetInt(order.TakerFee.Big())

	order.Price, _ = new(big.Float).Quo(takerTokenAmount, makerTokenAmount).Float64()
	order.FeeRate, _ = new(big.Float).Quo(takerFeeAmount, takerTokenAmount).Float64()

	if order.Cancelled {
		order.Status = StatusCancelled
	}

	if thresholdAmount.Cmp(big.NewInt(0)) <= 0 {
		order.Status = StatusFilled
	}
	if len(order.PoolID) == 0 {
		order.PoolID = DefaultSha3()
	}
}

// Save records the order in the database, defaulting to the specified status.
// Status should either be db.StatusOpen, or db.StatusUnfunded. If the order
// is filled based on order.TakerAssetAmountFilled + order.TakerAssetAmountCancelled
// the status will be recorded as db.StatusFilled regardless of the specified status.
func (order *Order) Save(db *gorm.DB, status int64, publisher channels.Publisher) *gorm.DB {
	if !order.Signature.Verify(order.Maker, order.Hash()) {
		scope := db.New()
		scope.AddError(errors.New("Failed to verify signature"))
		return scope
	}
	order.Populate()

	if order.Status == StatusOpen {
		order.Status = status
	}

	log.Printf("Attempting to save order %#x", order.Hash())



	updates := map[string]interface{}{
		"taker_asset_amount_filled":    order.TakerAssetAmountFilled,
		"maker_asset_remaining":        order.MakerAssetRemaining,
		"maker_fee_remaining":          order.MakerFeeRemaining,
		"status":                       order.Status,
	}

	updateScope := db.Model(Order{}).Where("order_hash = ?", order.OrderHash).Updates(updates)
	if updateScope.Error != nil {
		log.Printf(updateScope.Error.Error())
	}
	if updateScope.RowsAffected > 0 {
		if publisher != nil {
			publisher.Publish(string(order.Bytes()))
		}
		return updateScope
	}
	scope := db.Create(order)
	if scope.Error != nil {
		return scope
	}
	if publisher != nil {
		publisher.Publish(string(order.Bytes()))
	}
	return scope
}

type orderTracker struct {
	makerOrders []int
	takerOrders []int
}

func PopulateAssetMetadata(orders []Order, db *gorm.DB) {
	// TODO: See the code in asset_metadata for potentially replacing this with
	// subqueries.
	assetDataMap := make(map[string]*orderTracker)
	for i, order := range orders {
		if assetDataMap[fmt.Sprintf("%#x", order.MakerAssetData[:])] == nil {
			assetDataMap[fmt.Sprintf("%#x", order.MakerAssetData[:])] = &orderTracker{}
		}
		if assetDataMap[fmt.Sprintf("%#x", order.TakerAssetData[:])] == nil {
			assetDataMap[fmt.Sprintf("%#x", order.TakerAssetData[:])] = &orderTracker{}
		}
		assetDataMap[fmt.Sprintf("%#x", order.MakerAssetData[:])].makerOrders = append(assetDataMap[fmt.Sprintf("%#x", order.MakerAssetData[:])].makerOrders, i)
		assetDataMap[fmt.Sprintf("%#x", order.TakerAssetData[:])].takerOrders = append(assetDataMap[fmt.Sprintf("%#x", order.TakerAssetData[:])].takerOrders, i)
	}
	assetDataList := []*types.AssetData{}
	for _, v := range assetDataMap {
		if len(v.makerOrders) > 0 {
			assetDataList = append(assetDataList, &(orders[v.makerOrders[0]].MakerAssetData))
		} else {
			assetDataList = append(assetDataList, &(orders[v.takerOrders[0]].TakerAssetData))
		}
	}
	allAssetMetadata := []AssetMetadata{}
	db.Model(&AssetMetadata{}).Where("asset_data IN (?)", assetDataList).Find(&allAssetMetadata)
	for _, assetMetadata := range allAssetMetadata {
		for _, idx := range assetDataMap[fmt.Sprintf("%#x", assetMetadata.AssetData[:])].makerOrders {
			am := assetMetadata
			orders[idx].MakerAssetMetadata = &am
		}
		for _, idx := range assetDataMap[fmt.Sprintf("%#x", assetMetadata.AssetData[:])].takerOrders {
			am := assetMetadata
			orders[idx].TakerAssetMetadata = &am
		}
	}
}
