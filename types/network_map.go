package types

import (
	"database/sql/driver"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jinzhu/gorm"
	"errors"
	"log"
)
type NetworkAddressMap map[uint]*Address

type netAddPair struct {
	Net uint
	Address Address
}

func (nmap NetworkAddressMap) Value() (driver.Value, error) {
	items := []netAddPair{}
	for k, v := range nmap {
		items = append(items, netAddPair{k, *v})
	}
	bytes, err := rlp.EncodeToBytes(items)
	log.Printf("Saving nmap: %v, %#x, %v", items, bytes[:], err)
	bytes, err = rlp.EncodeToBytes(&items[0])
	log.Printf("nmap[0]: %v, %#x, %v", items[0], bytes[:], err)

	return rlp.EncodeToBytes(items)
}

func (nmap NetworkAddressMap) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		items := []netAddPair{}
		rlp.DecodeBytes(v, items)
		for _, item := range items {
			nmap[item.Net] = &item.Address
		}
		return nil
	default:
		return errors.New("NetworkAddressMap scanner src should be []byte")
	}
}


// GormDataType tells gorm what data type to use for the column.
func (nmap NetworkAddressMap) GormDataType(dialect gorm.Dialect) string {
	if dialect.GetName() == "postgres" {
		return "bytea"
	}
	return "blob"
}
