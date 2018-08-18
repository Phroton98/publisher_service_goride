package main

import (
	"log"
	"app.goride/model/order"
)

func main() {
	db, err := order.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	if !db.HasTable(&order.Order{}) {
		db.CreateTable(&order.Order{})
	}
	if !db.HasTable(&order.OrderLocation{}) {
		db.CreateTable(&order.OrderLocation{})
	}
	if !db.HasTable(&order.OrderFlag{}) {
		db.CreateTable(&order.OrderFlag{})
	}
}
