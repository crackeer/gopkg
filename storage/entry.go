package storage

import (
	"fmt"

	// MySQL Driver

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (

	// DriverMySQL
	DriverMySQL = "mysql"

	// DriverSQLite
	DriverSQLite = "sqlite"
)

// MySQLConfig
type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Database string
	Charset  string
}

// GetMySQLDB
//  @param dbConfig
//  @param gormConfig
//  @return *gorm.DB
//  @return *gorm.DB
//  @return error
func GetMySQLDB(dbConfig *MySQLConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Database, dbConfig.Charset)
	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}

	return gorm.Open(mysql.Open(dsn), gormConfig)
}

// GetSQliteDB
//  @param dbConfig
//  @param gormConfig
//  @return *gorm.DB
//  @return *gorm.DB
//  @return error
func GetSQliteDB(file string, gormConfig *gorm.Config) (*gorm.DB, error) {
	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}

	return gorm.Open(sqlite.Open(file), &gorm.Config{})
}
