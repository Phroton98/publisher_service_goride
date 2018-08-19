package order

import (
	"time"
	"app.goride/config"
	"github.com/jinzhu/gorm"
  	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Order struct {
	ID			uint `gorm:"AUTO_INCREMENT;unique_index;primary_key"`
	UserId		int `gorm:"NOT NULL"`
	DriverId	int
	Origin		string `gorm:"NOT NULL"`
	Destination	string `gorm:"NOT NULL"`
	Price		int `gorm:"NOT NULL;DEFAULT:0" `
	Status		string `gorm:"TYPE:VARCHAR(30);NOT NULL"`
	CreatedAt 	time.Time `gorm:"NOT NULL"`
	GoPay		bool
}

type OrderLocation struct {
	OrderID		uint `sql:"type:int REFERENCES orders(id)"`
	Order		Order `gorm:"foreignkey:OrderID"`
	OriginX		float64 `gorm:"NOT NULL;TYPE:DECIMAL(12,6)"`
	OriginY		float64 `gorm:"NOT NULL;TYPE:DECIMAL(12,6)"`
	DestX		float64 `gorm:"NOT NULL;TYPE:DECIMAL(12,6)"`
	DestY		float64 `gorm:"NOT NULL;TYPE:DECIMAL(12,6)"`
}

type OrderFlag struct {
	OrderID		uint  `sql:"type:int REFERENCES orders(id)"`
	Order		Order `gorm:"foreignkey:OrderID"`
	Flag		int
}

func (Order) TableName() string {
	return "orders"
}

func (OrderLocation) TableName() string {
	return "orders_location"
}

func (OrderFlag) TableName() string {
	return "orders_flag"
}

func ConnectDatabase() (db *gorm.DB, err error) {
	// var user, password, database string
	if config.HEROKU == "true" {
		db, err = gorm.Open("postgres", config.DATABASE_URL)
	} else {
		db, err = gorm.Open("postgres", "host=" + config.DB_HOST + " port=" + config.DB_PORT + " user=" + config.DB_USER + " dbname=" + config.DB_NAME + " password=" + config.DB_PASSWORD)
	}
	return db, err
}