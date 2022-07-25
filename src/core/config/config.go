package config

import (
	"log"
	"main/src/core/db"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	isDebug          bool
	webServerAddress string
	dbConfig         db.Config
}

func (c *Config) Init() error {
	log.Println("config: init")

	_ = godotenv.Load()

	err := c.loadDbConfig()
	if err != nil {
		return err
	}

	err = c.loadWebServerAddress()
	if err != nil {
		return err
	}

	err = c.loadIsDebug()
	if err != nil {
		return err
	}

	log.Println("config.init: successful")

	return nil
}

func (c *Config) GetDbConfig() db.Config {
	return c.dbConfig
}

func (c *Config) GetWebServerAddress() string {
	return c.webServerAddress
}

func (c *Config) IsDebug() bool {
	return c.isDebug
}

func (c *Config) loadDbConfig() error {
	var err error
	c.dbConfig = db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}

	maxOpenConnections := os.Getenv("DB_MAX_OPEN_CONNECTIONS")
	if len(maxOpenConnections) == 0 {
		c.dbConfig.MaxOpenConnections = 10
	} else {
		c.dbConfig.MaxOpenConnections, err = strconv.Atoi(maxOpenConnections)
		if err != nil {
			return err
		}
	}

	maxIdleConnections := os.Getenv("DB_MAX_IDLE_CONNECTIONS")
	if len(maxIdleConnections) == 0 {
		c.dbConfig.MaxIdleConnections = 10
	} else {
		c.dbConfig.MaxIdleConnections, err = strconv.Atoi(maxIdleConnections)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) loadWebServerAddress() error {
	c.webServerAddress = os.Getenv("WEB_SERVER_ADDR")
	if len(c.webServerAddress) == 0 {
		c.webServerAddress = "0.0.0.0:8080"
	}

	return nil
}

func (c *Config) loadIsDebug() error {
	isDebugVal := os.Getenv("DEBUG")
	c.isDebug = strings.Compare(isDebugVal, "true") == 0 || strings.Compare(isDebugVal, "1") == 0

	return nil
}

func NewConfig() *Config {
	return &Config{}
}
