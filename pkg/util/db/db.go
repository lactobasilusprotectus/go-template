package db

import (
	"fmt"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
	"time"
)

// DatabaseConnection is the connection to the database
type DatabaseConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

// NewDatabaseConnection constructs new DatabaseConnection
func NewDatabaseConnection(config config.DatabaseConfig) (*DatabaseConnection, error) {
	connStr := ConnStr(config)

	master, err := connect(connStr, config.MaxIdleConnections, config.MaxOpenConnections, config.Driver)
	if err != nil {
		return nil, err
	}

	slave, err := connect(connStr, config.MaxIdleConnections, config.MaxOpenConnections, config.Driver)
	if err != nil {
		return nil, err
	}

	return &DatabaseConnection{
		Master: master,
		Slave:  slave,
	}, nil
}

// ConnStr construct db connection string
func ConnStr(config config.DatabaseConfig) string {
	switch config.Driver {
	case "postgres":
		return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			config.Username, config.Password, config.Database, config.Host, config.Port)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
	case "mssql":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	case "sqlite":
		return fmt.Sprintf("%s", config.Database)
	}

	return "driver not supported"
}

// connect connects to database, given the configuration
func connect(address string, maxIdle, maxOpen int, driver string) (*gorm.DB, error) {
	//init variables
	var conn *gorm.DB
	var err error

	//decide which driver to use
	switch driver {
	case "postgres":
		conn, err = gorm.Open(postgres.Open(address), &gorm.Config{})
		break
	case "mysql":
		conn, err = gorm.Open(mysql.Open(address), &gorm.Config{})
		break
	case "mssql":
		conn, err = gorm.Open(sqlserver.Open(address), &gorm.Config{})
		break
	case "sqlite":
		conn, err = gorm.Open(sqlite.Open(address), &gorm.Config{})
		break
	default:
		log.Println("Error while connecting to database at", address, "driver not supported")
		return nil, fmt.Errorf("driver not supported")
	}

	if err != nil {
		log.Println("Error while connecting to database at", address, err.Error())
		return nil, err
	}

	//Set max idle and max open connections
	sqlDB, err := conn.DB()

	if err != nil {
		log.Println("Error while getting sql db", err.Error())
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdle)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpen)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return conn, nil
}

// CloseConnections closes database connections
func CloseConnections(dbConn *DatabaseConnection) {
	sqlDBMaster, err := dbConn.Master.DB()
	if err != nil {
		log.Println("Error while getting sql db", err.Error())
	}

	if err = sqlDBMaster.Close(); err != nil {
		log.Println("Error while closing sql db", err.Error())
	}

	sqlDBSlave, err := dbConn.Slave.DB()
	if err != nil {
		log.Println("Error while getting sql db", err.Error())
	}

	if err = sqlDBSlave.Close(); err != nil {
		log.Println("Error while closing sql db", err.Error())
	}
}
