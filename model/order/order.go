package order

import (
	"time"
	"errors"
	"net/http"
	"app.goride/config"
	"github.com/jinzhu/gorm"
)

// JSON for UserInformation Given
type OrderInformation struct {
	UserId int `json:"user_id" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Origin string `json:"origin" binding:"required"`
    OriginX float64 `json:"x" binding:"required"`
    OriginY float64 `json:"y" binding:"required"`
    Destination string `json:"destination" binding:"required"`
    DestX float64 `json:"dest_x" binding:"required"`
	DestY float64 `json:"dest_y" binding:"required"`
	Price int `json:"price" binding:"required"`
	GoPay *bool `json:"go_pay" binding:"exists"`
}

type CancelPayload struct {
	UserId int `json:"user_id" binding:"required"`
	TransactionId int `json:"transaction_id"`
}

type DeclinePayload struct {
	UserId int `json:"user_id" binding:"required"`
}

type UpdatePayload struct {
	Status string `json:"status" binding:"required"`
	UserId int `json:"user_id" binding:"required"`
	DriverId int `json:"driver_id" binding:"required"`
	TransactionId *int `json:"transaction_id,omitempty"`
}

func CreateOrder(u OrderInformation, totalDriver int) (*Order, error) {
	db, err := ConnectDatabase()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	order := Order{
		UserId: u.UserId,
		Origin: u.Origin,
		Destination: u.Destination,
		Price: u.Price,
		Status: config.QUEUEING,
		CreatedAt: time.Now(),
		GoPay: *u.GoPay,
	}
	db.Create(&order)
	orderLocation := OrderLocation{
		OrderID: order.ID,
		OriginX: u.OriginX,
		OriginY: u.OriginY,
		DestX: u.DestX,
		DestY: u.DestY,
	}
	orderFlag := OrderFlag{
		OrderID: order.ID,
		Flag: totalDriver,
	}
	db.Create(&orderLocation)
	db.Create(&orderFlag)
	return &order, nil
}

func GetOrder(id string) (*Order, error) {
	var order Order
	if db, err := ConnectDatabase(); err != nil {
		return nil, err
	} else {
		defer db.Close()
		if err = FindOrderById(db, id, &order); err == nil {
			return &order, nil
		} else {
			return nil, err
		}
	}
}

func FindOrderById(db *gorm.DB, id string, order *Order) (error) {
	db.Where("id = ?", id).First(order)
	if *order == (Order{}) {
		return errors.New("Order not found")
	} else {
		return nil
	}
}

func FindOrderFlagById(db *gorm.DB, id string, orderFlag *OrderFlag) (error) {
	db.Where("id = ?", id).First(orderFlag)
	if *orderFlag == (OrderFlag{}) {
		return errors.New("Order not found")
	} else {
		return nil
	}
}

func CancelOrder(id string, payload CancelPayload) (int, error) {
	if db, err := ConnectDatabase(); err != nil {
		return http.StatusInternalServerError, err
	} else {
		defer db.Close()
		var order Order
		if err = FindOrderById(db, id, &order); err != nil {
			return http.StatusNotFound, err
		} else if order.Status == config.CANCELLED || order.Status == config.FINISHED || order.Status == config.ACCEPTED {
			return http.StatusNotAcceptable, errors.New("Order cannot be canceled")
		} else if order.UserId != payload.UserId {
			return http.StatusNotAcceptable, errors.New("User not authenticated")
		}
		order.Status = config.CANCELLED
		db.Save(&order)
		return http.StatusNoContent, nil
	}
}

func AcceptOrder(id string, payload UpdatePayload) (int, error) {
	if db, err := ConnectDatabase(); err != nil {
		return http.StatusInternalServerError, err
	} else {
		defer db.Close()
		var order Order
		if err = FindOrderById(db, id, &order); err != nil {
			return http.StatusNotFound, err
		} else if order.Status != config.QUEUEING {
			return http.StatusNotAcceptable, errors.New("Order cannot be accepted!")
		}
		order.Status = config.ACCEPTED
		order.DriverId = payload.DriverId
		db.Save(&order)
		return http.StatusOK, nil
	}
}

func FinishOrder(id string, payload UpdatePayload) (int, error) {
	if db, err := ConnectDatabase(); err != nil {
		return http.StatusInternalServerError, err
	} else {
		defer db.Close()
		var order Order
		if err = FindOrderById(db, id, &order); err != nil {
			return http.StatusNotFound, err
		} else if payload.TransactionId == nil {
			return http.StatusNotAcceptable, errors.New("Transaction ID must be included!")
		} else if order.Status != config.ACCEPTED {
			return http.StatusNotAcceptable, errors.New("Order cannot be finished!")
		} 
		order.Status = config.FINISHED
		db.Save(&order)
		return http.StatusOK, nil
	}
}

func DeclineOrder(id string, payload DeclinePayload) (int, int, error) {
	if db, err := ConnectDatabase(); err != nil {
		return http.StatusInternalServerError, -1, err
	} else {
		defer db.Close()
		var orderFlag OrderFlag
		var order Order
		if err = FindOrderFlagById(db, id, &orderFlag); err != nil {
			return http.StatusNotFound, -1, errors.New("Order not found")
		}
		// Find related order
		db.Model(&orderFlag).Related(&order)
		if order.UserId != payload.UserId {
			return http.StatusNotFound, -1, errors.New("User doesn't match")
		}
		if (orderFlag.Flag <= 0) {
			orderFlag.Flag = 0
		} else {
			orderFlag.Flag = orderFlag.Flag - 1
		}
		db.Save(&orderFlag)
		return http.StatusOK, orderFlag.Flag, nil
	}
}