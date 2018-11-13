package types

import (
	"github.com/jinzhu/gorm"
)

type Pool interface {
	Filter(*gorm.DB) (*gorm.DB, error)
}
