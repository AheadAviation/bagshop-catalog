package api

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/AheadAviation/bagshop-catalog/db"
	"github.com/AheadAviation/bagshop-catalog/item"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	CreateItem(name, description string, price float32, count int) (string, error)
	Health() []Health
}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

type fixedService struct{}

func NewFixedService() Service {
	return &fixedService{}
}

func (s *fixedService) CreateItem(name, description string, price float32,
	count int) (string, error) {
	i := item.Item{}
	uid, _ := uuid.NewRandom()
	i.ID = uid.String()
	i.Name = name
	i.Description = description
	i.Price = price
	i.Count = count
	err := db.CreateItem(&i)
	return i.ID, err
}

func (s *fixedService) Health() []Health {
	var health []Health
	dbstatus := "OK"

	err := db.Ping()
	if err != nil {
		dbstatus = err.Error()
	}

	app := Health{"catalog", "OK", time.Now().String()}
	db := Health{"catalog-db", dbstatus, time.Now().String()}

	health = append(health, app)
	health = append(health, db)

	return health
}
