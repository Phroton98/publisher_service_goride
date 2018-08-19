package order

import "fmt"
import "bytes"
import "testing"
import "net/http"
import "io/ioutil"
import "encoding/json"

type RequestCreateOrder struct {
	UserId int `json:"user_id"`
	UserName string `json:"user_name"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Origin string `json:"origin"`
	Dest string `json:"destination"`
	DestX float64 `json:"dest_x"`
	DestY float64 `json:"dest_y"`
	Price int `json:"price"`
	GoPay bool `json:"go_pay"`
}

type RequestCancelOrder struct {
	UserId int `json:"user_id"`
	TransactionId int `json:"transaction_id"`
}

var host = "localhost"
var port = "8080"

func TestCreateOrder(t *testing.T) {
	urlPath := "http://" +  host + ":" + port + "/api/order"
	requestBody := RequestCreateOrder{
		UserId: 1,
		UserName: "compfest",
		X: 3.0,
		Y: 4.0001,
		Origin: "Grogol",
		Dest: "Puri Indah",
		DestX: 3.01,
		DestY: 4.01,
		Price: 10000,
		GoPay: true,
	}
	jsonValue, _ := json.Marshal(requestBody)
	request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	if response, err := client.Do(request); err != nil {
		t.Errorf("Testing Error. Get %s", err.Error())
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		defer response.Body.Close()
	}
}

func TestCancelOrder(t *testing.T) {
	urlPath := "http://" +  host + ":" + port + "/api/order/1"
	requestBody := RequestCancelOrder{1, -1}
	jsonValue, _ := json.Marshal(requestBody)
	request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	if response, err := client.Do(request); err != nil {
		t.Errorf("Testing Error. Get %s", err.Error())
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		defer response.Body.Close()
	}
}

func TestGetOrder(t *testing.T) {
	urlPath := "http://" +  host + ":" + port + "/api/order/1"
	if response, err := http.Get(urlPath); err != nil {
		t.Errorf("Testing Error. Get %s", err.Error())
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		defer response.Body.Close()
	}
}