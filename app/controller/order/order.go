package order

import (
    "fmt"
    "sort"
    "bytes"
    "strings"
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
    if err := c.ShouldBindJSON(&data); err == nil {
        if listDriverAvailable, errDB := getNearestDriver(data.OriginX, data.OriginY, MAX_DRIVER); errDB == nil {
            // Checking list Driver
            if len(listDriverAvailable) == 0 && config.ENVIRONMENT != "test" {
                c.JSON(http.StatusNotFound, gin.H{"error": "drivers not found"})
                return
            }
            if orderData, err := order.CreateOrder(data, MAX_DRIVER); err == nil {
                // TO DO
                // POST To Gopay
                var id int
                if id, err = CreateTransaction(orderData, data.UserName); err == nil {
                    // POST to Subscriber
                    distance := helper.GetDistance(data.OriginX, data.OriginY, data.DestX, data.DestY)
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
                    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                } 
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message":err.Error()})
    }
}

// Tested
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
    intId, _ := strconv.Atoi(orderID)
    var data order.CancelPayload
    if err := c.ShouldBindJSON(&data); err == nil {
        var status int
        if status, err = order.CancelOrder(orderID, data); err == nil {
            // TO DO
            // Send request to Subscriber
            if err = SendInvalidate(intId); err == nil {
                c.JSON(status, gin.H{"message": "order cancelled"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
            }
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
        intId, _ := strconv.Atoi(id)
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
    if status, err := order.FinishOrder(id, data); err == nil {
        // To Do
        // Post to Go Pay Wallet
        path := "/api/transaction/" + strconv.Itoa(*data.TransactionId)
        urlPath := CreateGopayPath(path)
        // Create data request
        dataRequest := url.Values{}
        dataRequest.Add("finished", "true")
        // Creat client
        client := &http.Client{}
        request, _ := http.NewRequest("POST", urlPath, strings.NewReader(dataRequest.Encode()))
        if resp, errResp := client.Do(request); errResp == nil {
            defer resp.Body.Close()
            // Send response
            c.JSON(status, gin.H{"message": "order finished"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error_message": errResp.Error()})
        }
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