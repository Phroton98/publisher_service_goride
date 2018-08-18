package migration

import (
	"log"
	"app.goride/model/driver"
	"app.goride/model/order"
)

func CreateAllTable() {
	db, err := order.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	db.DropTableIfExists(&order.OrderFlag{}, &order.OrderLocation{}, &order.Order{}, &driver.DriverLocation{})
	
	if !db.HasTable(&order.Order{}) {
		db.CreateTable(&order.Order{})
	}
	if !db.HasTable(&order.OrderLocation{}) {
		db.CreateTable(&order.OrderLocation{})
	}
	if !db.HasTable(&order.OrderFlag{}) {
		db.CreateTable(&order.OrderFlag{})
	}
	if !db.HasTable(&driver.DriverLocation{}) {
		db.CreateTable(&driver.DriverLocation{})
	}
}