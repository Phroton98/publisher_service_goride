package test

import "testing"
import "fmt"
import "io/ioutil"
import "bytes"
import "encoding/json"
import "net/http"

type Request struct {
	X int `json:"x"`
	Y int `json:"y"`
	Available bool `json:"available"`
}

var host = "localhost"
var port = "8080"

func TestPostLocation(t *testing.T) {
	urlPath := "http://" +  host + ":" + port + "/api/location/999"
	requestBody := Request{
		X: 4.0,
		Y: 5.0,
		Available: true,
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

func TestGetLocation(t *testing.T) {
	urlPath := "http://" +  host + ":" + port + "/api/location/999"
	if response, err := http.Get(urlPath); err != nil {
		t.Errorf("Testing Error. Get %s", err.Error())
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		defer response.Body.Close()
	}
}