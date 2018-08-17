package order

import (
    "fmt"
    "net/http"
    "sort"
    "encoding/json"
    "app.goride/model/driver"
    "app.goride/model/order"
    "github.com/gin-gonic/gin"
)

// Constant
const THRESHOLD = 1000 // In Metres
const MAX_DRIVER = 30

func CreateOrder(c *gin.Context) {
    var data order.OrderInformation
    var id int
    if err := c.ShouldBindJSON(&data); err == nil {
        listDriverAvailable := getNearestDriver(data.OriginX, data.OriginY, MAX_DRIVER)
        // if len(listDriverAvailable) == 0 {
        //     c.JSON(http.StatusNotFound, gin.H{"error": "drivers not found"})
        //     return
        // }
        driversJSON, _ := json.Marshal(listDriverAvailable)
        fmt.Println(driversJSON)
        if id, err = order.CreateOrder(data, MAX_DRIVER); err == nil {
            // TO DO
            // POST to Subscriber
            // POST to GoPay
            c.JSON(http.StatusCreated, gin.H{
                "message": "order created",
                "order_id": id,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message":err.Error()})
    }
}

func GetOrder(c *gin.Context) {
    orderID := c.Param("id")
    if order, err := order.GetOrder(orderID); err == nil {
        c.JSON(http.StatusOK, &order)
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}

func CancelOrder(c *gin.Context) {
    orderID := c.Param("id")
    var data order.CancelPayload
    if err := c.ShouldBindJSON(&data); err == nil {
        var status int
        if status, err = order.CancelOrder(orderID, data); err == nil {
            c.JSON(status, gin.H{"message": "order cancelled"})
            // TO DO
            // Send request to Subscriber
        } else {
            c.JSON(status, gin.H{"error_message": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}

func UpdateOrder(c *gin.Context) {
    orderID := c.Param("id")
    var data order.UpdatePayload
    if err := c.ShouldBindJSON(&data); err == nil {
        if data.Status == "accept" {
            AcceptOrder(orderID, data, c)
        } else if data.Status == "finish" {
            FinishOrder(orderID, data, c)
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error_message": "Status value must accepted or finished!"})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}

func AcceptOrder(id string, data order.UpdatePayload, c *gin.Context) {
    if status, err := order.AcceptOrder(id, data); err == nil {
        c.JSON(status, gin.H{"message": "order accepted"})
        // To Do
        // Post to subscriber
    } else {
        c.JSON(status, gin.H{"error_message": err.Error()})
    }
}

func FinishOrder(id string, data order.UpdatePayload, c *gin.Context) {
    if status, err := order.AcceptOrder(id, data); err == nil {
        c.JSON(status, gin.H{"message": "order finished"})
        // To Do
        // Post to Go Pay Wallet
    } else {
        c.JSON(status, gin.H{"error_message": err.Error()})
    }
}

func DeclineOrder(c *gin.Context) {
    orderID := c.Param("id")
    var data order.DeclinePayload
    if err := c.ShouldBindJSON(&data); err == nil {
        var flag, status int
        // Return status and flag
        if status, flag, err = order.DeclineOrder(orderID, data); err != nil {
            if flag == 0 {
                // TO DO
                // Finding Drivers AGAIN
            }
            c.JSON(status, gin.H{"message": "order declined"})
        } else {
            c.JSON(status, gin.H{"error": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}

func getNearestDriver(clientX float64, clientY float64, maxDriver int) []driver.DriverInformation {
    listDriver := driver.GetDriverAround(THRESHOLD, clientX, clientY)
    // Sort
    sort.Slice(listDriver, func(i int, j int) bool {
        return listDriver[i].Distance < listDriver[j].Distance
    })
    if len(listDriver) > maxDriver {
        listDriver = listDriver[:maxDriver]
    }
    return listDriver
}