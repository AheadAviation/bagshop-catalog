package api

import (
	"errors"
	"time"

	"github.com/AheadAviation/bagshop-catalog/db"
)

var ErrNotFound = errors.New("not found")

type Service interface {
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
