package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/namsral/flag"
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
	MySQLc *sql.DB
}

func (m *MySQL) Init() error {
	var dsn string
	if username != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, addr, db)
	} else {
		dsn = fmt.Sprintf("@tcp(%s)/%s", addr, db)
	}

	var err error
	m.MySQLc, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQL) Ping() error {
	return m.MySQLc.Ping()
}
