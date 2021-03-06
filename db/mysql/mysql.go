package mysql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"

	"github.com/AheadAviation/bagshop-catalog/item"
)

var (
	username string
	password string
	addr     string
	db       = "catalog"
)

func init() {
	flag.StringVar(&username, "mysql-username", "", "MySQL DB Username")
	flag.StringVar(&password, "mysql-password", "", "MySQL DB Password")
	flag.StringVar(&addr, "mysql-addr", "0.0.0.0:3306", "MySQL Host Address and Port")
}

type MySQL struct {
	MySQLc *gorm.DB
}

func (m *MySQL) Init() error {
	var dsn string
	if username != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			username, password, addr, db)
	} else {
		dsn = fmt.Sprintf("@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			addr, db)
	}

	var err error
	// m.MySQLc, err = sql.Open("mysql", dsn)
	m.MySQLc, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	m.MySQLc.AutoMigrate(&item.Item{})
	return m.seedData()
}

func (m *MySQL) CreateItem(i *item.Item) error {
	return m.MySQLc.Create(i).Error
}

func (m *MySQL) GetItems() ([]item.Item, error) {
	its := make([]item.Item, 0)
	r := m.MySQLc.Find(&its)
	return its, r.Error
}

func (m *MySQL) Ping() error {
	return m.MySQLc.DB().Ping()
}

func (m *MySQL) seedData() error {
	its := make([]item.Item, 0)
	m.MySQLc.Find(&its)
	if len(its) == 0 {
		sd, err := ioutil.ReadFile("/seed-data.json")
		if err != nil {
			return err
		}

		json.Unmarshal(sd, &its)

		for i := range its {
			err := m.CreateItem(&its[i])
			if err != nil {
				return err
			}
		}
		log.Printf("Seeded %d items into the database", len(its))
	}
	return nil
}
