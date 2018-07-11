package funds

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	dbModule "github.com/notegio/openrelay/db"
	"bytes"
)

type CancellationLookup interface {
	GetCancelled(order *types.Order) (bool, error)
}

type dbCancellationLookup struct {
	db *gorm.DB
}

func (lookup *dbCancellationLookup) GetCancelled(order *types.Order) (bool, error) {
	if(order.Cancelled){
		// If it was cancelled earlier, it's still going to be cancelled. We're
		// mostly interested in things that have been cancelled since the last
		// check.
		return true, nil
	}
	cancellation := &dbModule.Cancellation{}
	if err := lookup.db.Model(&dbModule.Cancellation{}).Where("maker = ? AND sender = ?", order.Maker, order.SenderAddress).First(cancellation).Error; err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return true, err
	}
	if cancellation.Epoch == nil {
		return false, nil
	}
	return bytes.Compare(order.Salt[:], cancellation.Epoch[:]) < 0, nil
}

func NewDBCancellationLookup(db *gorm.DB) (CancellationLookup) {
	return &dbCancellationLookup{db}
}

type MockCancellationLookup struct {
	cancelled bool
}

func (filled *MockCancellationLookup) GetCancelled(order *types.Order) (bool, error) {
	return filled.cancelled, nil
}

func NewMockCancellationLookup(cancelled bool) CancellationLookup {
	return &MockCancellationLookup{cancelled}
}
