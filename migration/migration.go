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
	if !db.HasTable(&order.Order{}) {
		db.DropTable(&order.Order{})
		db.CreateTable(&order.Order{})
	}
	if !db.HasTable(&order.OrderLocation{}) {
		db.DropTable(&order.OrderLocation{})
		db.CreateTable(&order.OrderLocation{})
	}
	if !db.HasTable(&order.OrderFlag{}) {
		db.DropTable(&order.OrderFlag{})
		db.CreateTable(&order.OrderFlag{})
	}
	if !db.HasTable(&driver.DriverLocation{}) {
		db.DropTable(&driver.DriverLocation{})
		db.CreateTable(&driver.DriverLocation{})
	}
}