package driver

import (
	"app.goride/config"
	"github.com/jinzhu/gorm"
  	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DriverLocation struct {
    ID int `json:"id" gorm:"unique_index"`
    Token string `json:"token" gorm:"TYPE:VARCHAR(100);NOT NULL"`
    X float64 `json:"x" gorm:"NOT NULL"`
    Y float64 `json:"y" gorm:"NOT NULL"`
    Available *bool `json:"available" gorm:"NOT NULL"`
    Timestamp int64 `json:"timestamp" gorm:"NOT NULL"`
}

func (DriverLocation) TableName() string {
	return "locations"
}

func ConnectDatabase() (db *gorm.DB, err error) {
	// var user, password, database string
	db, err = gorm.Open("postgres", config.DATABASE_URL)
	return db, err
}