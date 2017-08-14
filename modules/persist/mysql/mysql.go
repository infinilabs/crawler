package mysql

import (
	"github.com/jinzhu/gorm"
)

// MySQLConfig defines mysql related config, currently only provide connection, eg: root:password@tcp(127.0.0.1:3306)/gopa?charset=utf8&parseTime=true&loc=Local
type MySQLConfig struct {
	Connection string `config:"connection"` //TODO, move to structured config
}

// GetInstance return mysql instance for further usage
// before use MySQL, you should create database first
//DROP DATABASE gopa;
//CREATE DATABASE IF NOT EXISTS gopa DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
func GetInstance(cfg *MySQLConfig) *gorm.DB {

	var err error
	db, err := gorm.Open("mysql", cfg.Connection)
	if err != nil {
		panic("failed to connect database")
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	return db
}
