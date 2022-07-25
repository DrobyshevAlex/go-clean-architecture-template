package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"log"
	"time"
)

type Db struct {
	config     Config
	connection *sql.DB
}

func (c *Db) Init(config Config) {
	log.Print("db: init")

	c.config = config

	log.Println("db.init: success.")
}

func (c *Db) Run() {
	log.Println("db.run: connect to server")

	var connectionString string
	if len(c.config.Password) == 0 {
		connectionString = fmt.Sprintf(
			"%s@tcp(%s:%s)/%s?parseTime=true",
			c.config.Username,
			c.config.Host,
			c.config.Port,
			c.config.Database,
		)
	} else {
		connectionString = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.config.Username,
			c.config.Password,
			c.config.Host,
			c.config.Port,
			c.config.Database,
		)
	}

	var err error
	c.connection, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	c.connection.SetConnMaxLifetime(time.Minute * 3)
	c.connection.SetMaxOpenConns(c.config.MaxOpenConnections)
	c.connection.SetMaxIdleConns(c.config.MaxIdleConnections)

	err = c.connection.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("db.run: connect success.")
}

func (c *Db) GetClient() *sql.DB {
	return c.connection
}

func (c *Db) IsReady() bool {
	err := c.connection.Ping()
	if err != nil {
		return false
	}

	return true
}

func NewDb() *Db {
	return &Db{}
}
