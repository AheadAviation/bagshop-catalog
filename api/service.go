package catalog

import (
	"errors"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	List(tags []string, order string, pageNum, pageSize int) ([]Bag, error)
	Count(tags []string) (int, error)
	Get(id string) (Bag, error)
	Tags() ([]string, error)
	Health() []Health
}

type Middleware func(Service) Service

type Bag struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	ImageURL    []string `json:"imageUrl" db:"-"`
	ImageURL_1  string   `json:"-" db:"image_url_1"`
	ImageURL_2  string   `json:"-" db:"image_url_2"`
	Price       float32  `json:"price" db:"price"`
	Count       int      `json:"count" db:"count"`
	Tags        []string `json:"tag" db:"-"`
	TagString   string   `json:"-" db:"tag_name"`
}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

var ErrNotFound = errors.New("not found")

var ErrDBConnection = errors.New("database connection error")

var baseQuery = "SELECT bag.bag_id AS id, bag.name, bag.description, bag.price, " +
	"bag.count, bag.image_url_1, bag.image_url_2, GROUP_CONCAT(tag.name) " +
	"AS tag_name FROM bag JOIN bag_tag ON bag.bag_id=bag_tag.bag_id JOIN " +
	"tag ON bag_tag.tag_id=tag.tag_id"

func NewCatalogService(db *sqlx.DB, logger log.Logger) Service {
	return &catalogService{
		db:     db,
		logger: logger,
	}
}

type catalogService struct {
	db     *sqlx.DB
	logger log.Logger
}

func (s *catalogService) List(tags []string, order string, pageNum, pageSize int) ([]Bag, error) {
	var bags []Bag
	query := baseQuery

	var args []interface{}

	for i, t := range tags {
		if i == 0 {
			query += " WHERE tag.name=?"
			args = append(args, t)
		} else {
			query += " OR tag.name=?"
			args = append(args, t)
		}
	}

	query += " GROUP BY id"

	if order != "" {
		query += " ORDER BY ?"
		args = append(args, order)
	}

	query += ";"

	err := s.db.Select(&bags, query, args...)
	if err != nil {
		s.logger.Log("database error", err)
		return []Bag{}, ErrDBConnection
	}
	for i, s := range bags {
		bags[i].ImageURL = []string{s.ImageURL_1, s.ImageURL_2}
		bags[i].Tags = strings.Split(s.TagString, ",")
	}

	bags = cut(bags, pageNum, pageSize)

	return bags, nil
}

func (s *catalogService) Count(tags []string) (int, error) {
	query := "SELECT COUNT(DISTINCT bag.bag_id) FROM bag " +
		"JOIN bag_tag ON bag.bag_id=bag_tag.bag_id " +
		"JOIN tag ON bag_tag.tag_id=tag.tag_id"

	var args []interface{}

	for i, t := range tags {
		if i == 0 {
			query += " WHERE tag.name=?"
			args = append(args, t)
		} else {
			query += " OR tag.name=?"
			args = append(args, t)
		}
	}

	query += ";"

	sel, err := s.db.Prepare(query)
	if err != nil {
		s.logger.Log("database error", err)
		return 0, ErrDBConnection
	}
	defer sel.Close()

	var count int
	err = sel.QueryRow(args...).Scan(&count)

	if err != nil {
		s.logger.Log("database error", err)
		return 0, ErrDBConnection
	}

	return count, nil
}

func cut(bags []Bag, pageNum, pageSize int) []Bag {
	if pageNum == 0 || pageSize == 0 {
		return []Bag{}
	}
	start := (pageNum * pageSize) - pageSize
	if start > len(bags) {
		return []Bag{}
	}
	end := (pageNum * pageSize)
	if end > len(bags) {
		end = len(bags)
	}
	return bags[start:end]
}
