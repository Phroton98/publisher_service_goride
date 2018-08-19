package order

import (
    "fmt"
    "sort"
    "bytes"
    "strconv"
    "net/url"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "app.goride/config"
    "app.goride/app/helper"
    "app.goride/model/driver"
    "app.goride/model/order"
    "github.com/gin-gonic/gin"
)

// Constant
const THRESHOLD = 1000 // In Metres
const MAX_DRIVER = 30

func CreateGopayPath(path string) (string) {
    u, _ := url.ParseRequestURI(config.API_GOPAY)
    u.Path = path
    return u.String()
}

func CreateSubsPath(path string) (string) {
    u, _ := url.ParseRequestURI(config.API_SUBSCRIBER)
    u.Path = path
    return u.String()
}

func CreateTransaction(orderData *order.Order, username string) (int, error) {
    // Create path and body
    urlPath := CreateGopayPath("/api/transactions/")
    requestBody := order.CreateTransactionBody(username, orderData)
    jsonValue, _ := json.Marshal(requestBody)
    // Create Request
    request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    // Create client
    client := &http.Client{}
    if response, err := client.Do(request); err == nil {
        defer response.Body.Close()
        body, _ := ioutil.ReadAll(response.Body)
        // Response struct
        var responseStruct order.RespCreateTransaction
        json.Unmarshal(body, &responseStruct)
        return responseStruct.ID, nil
    } else {
        return -1, err
    }
}

func CancelTransaction(id string) {
    urlPath := CreateGopayPath("/api/transactions/" + id)
    request, _ := http.NewRequest("DELETE", urlPath, bytes.NewBuffer([]byte("")))
    client := &http.Client{}
    response, _ := client.Do(request)
    body, _ := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))
    response.Body.Close()
}

func FinishTransaction(id string) {
    urlPath := CreateGopayPath("/api/transactions/" + id)
    jsonValue := []byte(`{"finished":true}`)
    // Create Request
    request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValue))
    client := &http.Client{}
    response, _ := client.Do(request)
    body, _ := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))
    response.Body.Close()
}

func SendOrderSubscriber(orderData *order.Order, listDriver []driver.DriverInformation, distance int, transID int) (error) {
    urlPath := CreateSubsPath("/order")
    requestBody := order.CreateRequestSubs(orderData, listDriver, distance, transID)
    jsonValue, _ := json.Marshal(requestBody)
    fmt.Println(string(jsonValue))
    // Create Request
    request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    // Create Client
    client := &http.Client{}
    if response, err := client.Do(request); err == nil {
        defer response.Body.Close()
        body, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(body))
        return nil
    } else {
        return err
    }
}

func SendInvalidate(id int) (error) {
    // Create path and body
    urlPath := CreateSubsPath("/invalidate")
    requestBody := order.CreateRequestInvalid(id)
    jsonValue, _ := json.Marshal(requestBody)
    // Create Request
    request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    // Create Client
    client := &http.Client{}
    if response, err := client.Do(request); err == nil {
        defer response.Body.Close()
        body, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(body))
        return nil
    } else {
        return err
    }
}

func CreateOrder(c *gin.Context) {
    var data order.OrderInformation
    err := c.ShouldBindJSON(&data)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error_message":err.Error()})
        return 
    }
    // Get nearest driver
    listDriverAvailable, errDB := getNearestDriver(data.OriginX, data.OriginY, MAX_DRIVER)
    if errDB != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // Check if there is driver available
    if len(listDriverAvailable) == 0 && config.ENVIRONMENT != "test" {
        c.JSON(http.StatusNotFound, gin.H{"error": "drivers not found"})
        return
    }
    // Get distance
    distance := helper.GetDistance(data.OriginX, data.OriginY, data.DestX, data.DestY)
    if orderData, err := order.CreateOrder(data, MAX_DRIVER); err == nil {
        // POST To Gopay
        var id = -1
        // Check if pay with GoPay
        if (*data.GoPay == true) {
            id, err = CreateTransaction(orderData, data.UserName)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return 
            } 
        } 
        // POST to Subscriber
        if err = SendOrderSubscriber(orderData, listDriverAvailable, distance, id); err == nil {
            // Response OK
            c.JSON(http.StatusCreated, gin.H{
                "message": "order created",
                "order_id": int(orderData.ID),
                "transaction_id": id,
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
    // Get request body
    orderID := c.Param("id")
    intId, _ := strconv.Atoi(orderID)
    var data order.CancelPayload
    err := c.ShouldBindJSON(&data)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
        return
    }
    // Cancel Order
    var status int
    status, err = order.CancelOrder(orderID, data)
    if err != nil {
        c.JSON(status, gin.H{"error_message": err.Error()})
        return
    }
    // Cancel Transaction
    if data.TransactionId > 0 {
        CancelTransaction(strconv.Itoa(data.TransactionId))
    }
    // Send request to Subscriber
    if err = SendInvalidate(intId); err == nil {
        c.JSON(status, gin.H{"message": "order canceled"})
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
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
    if errDriver := driver.ChangeAvailable(data.DriverId, false); errDriver != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error_message": errDriver.Error()})
        return 
    }
    // Accepting Order
    if status, err := order.AcceptOrder(id, data); err == nil {
        intId, _ := strconv.Atoi(id)
        // Post to subscriber
        if err = SendInvalidate(intId); err == nil {
            c.JSON(status, gin.H{"message": "order accepted"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
        }
    } else {
        c.JSON(status, gin.H{"error_message": err.Error()})
    }
}

func FinishOrder(id string, data order.UpdatePayload, c *gin.Context) {
    if errDriver := driver.ChangeAvailable(data.DriverId, false); errDriver != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error_message": errDriver.Error()})
        return 
    }
    status, err := order.FinishOrder(id, data)
    if err != nil {
        c.JSON(status, gin.H{"error_message": err.Error()})
    }
    // Post to Go Pay Wallet if transaction id > 0
    if (*data.TransactionId > 0) {
        FinishTransaction(strconv.Itoa(*data.TransactionId))
    }
    c.JSON(status, gin.H{"message": "order finished"})
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

func getNearestDriver(clientX float64, clientY float64, maxDriver int) ([]driver.DriverInformation, error) {
    if listDriver, err := driver.GetDriverAround(THRESHOLD, clientX, clientY); err == nil {
        // Sort
        sort.Slice(listDriver, func(i int, j int) bool {
            return listDriver[i].Distance < listDriver[j].Distance
        })
        if len(listDriver) > maxDriver {
            listDriver = listDriver[:maxDriver]
        }
        return listDriver, nil
    } else {
        return nil, err
    }
  
}