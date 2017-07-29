package mysql

import (
	"github.com/jinzhu/gorm"
)

type MySQLConfig struct {
	Connection string `config:"connection"` //TODO, move to structured config
}

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
