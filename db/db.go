package db

import (
	"errors"
	"fmt"

	"github.com/namsral/flag"

	"github.com/AheadAviation/bagshop-catalog/item"
)

var (
	dbEngine              string
	DefaultDB             Database
	DBTypes               = map[string]Database{}
	ErrNoDatabaseFound    = "No database with name %v registered"
	ErrNoDatabaseSelected = errors.New("No DB selected")
)

type Database interface {
	Init() error
	CreateItem(*item.Item) error
	GetItems() ([]item.Item, error)
	Ping() error
}

func init() {
	flag.StringVar(&dbEngine, "database_engine", "mysql", "Database Engine Name (Only MySQL Currently Supported)")
}

func Init() error {
	if dbEngine == "" {
		return ErrNoDatabaseSelected
	}

	err := Set()
	if err != nil {
		return err
	}
	return DefaultDB.Init()
}

func Set() error {
	if v, ok := DBTypes[dbEngine]; ok {
		DefaultDB = v
		return nil
	}
	return fmt.Errorf(ErrNoDatabaseFound, dbEngine)
}

func Register(name string, db Database) {
	DBTypes[name] = db
}

func CreateItem(i *item.Item) error {
	return DefaultDB.CreateItem(i)
}

func GetItems() ([]item.Item, error) {
	return DefaultDB.GetItems()
}

func Ping() error {
	return DefaultDB.Ping()
}
