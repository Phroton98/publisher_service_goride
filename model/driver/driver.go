// Package model user
package driver

import (
    "time"
    "strconv"
    "encoding/json"
    "app.goride/app/helper"
    "github.com/go-redis/redis"
)

// JSON for UpdateLocation
type Location struct {
    X float64 `json:"x" binding:"required"`
    Y float64 `json:"y" binding:"required"`
    Available bool `json:"available" binding:"required"`
    Token string `json:"token" binding:"required"` 
}

type DriverLocation struct {
    ID int `json:"id"`
    Token string `json:"token" binding:"required"`
    X float64 `json:"x"`
    Y float64 `json:"y"`
    Available bool `json:"available"`
    Timestamp int64 `json:"timestamp"`
}

type DriverInformation struct {
    ID int `json:"id"`
    Token string `json:"token"`
    Distance int `json:"distance"`
}

func CreateDriverInformation(loc DriverLocation, distance int) DriverInformation {
    return DriverInformation{
        ID: loc.ID, 
        Token: loc.Token,
        Distance: distance,
    }
}   

func CreateDriverLocation(id int, data Location) DriverLocation {
    return DriverLocation{
        ID: id, 
        Token: data.Token,
        X: data.X, 
        Y: data.Y,
        Available: data.Available, 
        Timestamp: time.Now().Unix(),
    }
}

func createClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", 
        DB:       0,  
    })
}

func SetLocation(data *DriverLocation) (err error) {
    client := createClient()
    if b, err := json.Marshal(data); err != nil {
        return err
    } else if err = client.Set("driver:" + strconv.Itoa(data.ID), string(b), 0).Err(); err != nil {
        return err
    }
    return nil
}

func GetLocation(id string) (*DriverLocation, error) {
    client := createClient()
    if val, err := client.Get("driver:" + id).Result(); err == nil {
        var response DriverLocation
        if err = json.Unmarshal([]byte(val), &response); err == nil {
            return &response, nil
        } else {
            return nil, err
        }
    } else {
        return nil, err
    }
}

func GetDriverAround(threshold int, clientX float64, clientY float64) []DriverInformation {
    // Create client
    client := createClient()
    // Initialize drivers
    drivers := []DriverInformation{}
    // Set cursor
    var cursor uint64
    var err error
    for {
        var keys []string
        if keys, cursor, err = client.Scan(cursor, "driver:*", 50).Result(); err != nil {
            return []DriverInformation{}
        }
        if len(keys) > 0 {
            for _, key := range keys {
                // Define value dan data
                var val string
                var data DriverLocation
                // Get value from redis key
                val, err = client.Get(key).Result()
                // Change to struct
                if err = json.Unmarshal([]byte(val), &data); err != nil {
                    return []DriverInformation{}
                }
                if !data.Available {
                    continue
                }
                distance := helper.GetDistance(clientX, clientY, data.X, data.Y)
                if int(distance) <= threshold {
                    driver := CreateDriverInformation(data, int(distance))
                    drivers = append(drivers, driver)
                }
            }
        }
        if cursor == 0 {
            break
        }
    }
    return drivers
}