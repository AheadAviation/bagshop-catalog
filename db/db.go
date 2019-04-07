package db

import (
	"github.com/aheadaviation/bagshop-catalog/item"
)

type Database interface {
	GetItems() ([]item.Item, error)
}
