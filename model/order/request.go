package order

import (
	"app.goride/model/driver"
)

type RequestSubs struct {
	OrderID int `json:"OrderId"`
	Origin string `json:"Origin"`
	Destination string `json:"Destination"`
	Distance int `json:"DestinationDistance"`
	TransactionID int `json:"TransactionID"`
	DriverData []driver.DriverInformation `json:"DriverData"`
}

type RequestInvalid struct {
	OrderID int `json:"OrderId"`
}

type RequestGopay struct {
	User string `json:"user"`
	Balance int `json:"changed_balance"`
	Desc string `json:"description"`
}

func CreateTransactionBody(username string, order *Order) (RequestGopay) {
	return RequestGopay{
		User: username,
		Balance: order.Price * -1,
		Desc: "destination to " + order.Destination,
	}
}

func CreateRequestSubs(order *Order, listDriver []driver.DriverInformation, distance int, transID int) (RequestSubs) {
	return RequestSubs{
		OrderID: int(order.ID),
		Origin: order.Origin,
		Destination: order.Destination,
		Distance: distance,
		TransactionID: transID,
		DriverData: listDriver,
	}
}

func CreateRequestInvalid(id int) (RequestInvalid) {
	return RequestInvalid{OrderID: id}
}