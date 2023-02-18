package db

import (
	"fmt"
	"github.com/paccolamano/go-backend-common/env"
	"gorm.io/gorm/schema"
	"strconv"
	"time"

	"github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	Conn() *gorm.DB
	HealthCheck() healthcheck.Check
}

type database struct {
	connection *gorm.DB
	user       string
	password   string
	host       string
	port       string
	schema     string
}

func (db *database) Conn() *gorm.DB {
	return db.connection
}

func (db *database) HealthCheck() healthcheck.Check {
	sqlDb, err := db.connection.DB()
	if err != nil {
		log.Errorf("unable to create a healthcheck for database connection")
	}
	return healthcheck.DatabasePingCheck(sqlDb, time.Second)
}

func createDatabase() Database {
	var db *gorm.DB = nil

	// TODO throw error if env is not set
	retry, _ := strconv.Atoi(env.GetEnv("MYSQL_CONN_RETRY", ""))
	dbUser := env.GetEnv("MYSQL_USER", "")
	dbPassword := env.GetEnv("MYSQL_PASSWORD", "")
	dbHost := env.GetEnv("MYSQL_HOST", "")
	dbPort := env.GetEnv("MYSQL_PORT", "")
	dbSchema := env.GetEnv("MYSQL_DATABASE", "")

	var connected = false
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbSchema)
	for i := 0; !connected && i < retry; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			NamingStrategy: schema.NamingStrategy{
				NoLowerCase: true,
			},
		})
		if err != nil {
			log.Errorf("cannot connect to db '%s'. attempt number %d of %d failed.", dbSchema, i+1, retry)
			time.Sleep(5 * time.Second)
			continue
		}
		connected = true
	}

	if !connected {
		log.Error("cannot connect to database. exiting...")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("connection pool cannot be set up due to: %s", err)
	}
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &database{
		connection: db,
	}
}
