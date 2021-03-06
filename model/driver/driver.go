// Package model user
package driver

import (
    // "fmt"
    "time"
    "errors"
    "strconv"
    "app.goride/app/helper"
)

// JSON for UpdateLocation
type Location struct {
    X float64 `json:"x" binding:"required"`
    Y float64 `json:"y" binding:"required"`
    Available bool `json:"available" binding:"required"`
}

type DriverInformation struct {
    ID int `json:"DriverID"`
    Distance int `json:"OriginDistance"`
}

func CreateDriverInformation(loc DriverLocation, distance int) DriverInformation {
    return DriverInformation{
        ID: loc.ID, 
        Distance: distance,
    }
}   

func CreateDriverLocation(data Location, strID string) (DriverLocation) {
    id, _ := strconv.Atoi(strID)
    return DriverLocation{
        ID: id,
        X: data.X,
        Y: data.Y,
        Available: &data.Available,
        Timestamp: time.Now().Unix(),
    }
}

func SetLocation(data Location, id string) (error) {
    if db, err := ConnectDatabase(); err == nil {
        defer db.Close()
        var driver DriverLocation
        db.Where("id = ?", id).First(&driver)
        // Check if empty
        if driver == (DriverLocation{}) {
            driver = CreateDriverLocation(data, id)
            db.Create(&driver)
        } else {
            driver.X = data.X
            driver.Y = data.Y
            driver.Timestamp = time.Now().Unix()
            db.Save(&driver)
        }
    } else {
        return err
    }
    return nil
}

func GetLocation(id string) (*DriverLocation, error) {
    if db, err := ConnectDatabase(); err == nil {
        defer db.Close()
        var driver DriverLocation
        db.Where("id = ?", id).First(&driver)
        if driver != (DriverLocation{}) {
            return &driver, nil
        } else {
            return nil, errors.New("Driver not found!")
        }
    } else {
        return nil, err
    }
}

func GetDriverAround(threshold int, clientX float64, clientY float64) ([]DriverInformation, error) {
    if db, err := ConnectDatabase(); err == nil {
        defer db.Close()
        listDriverLoc := []DriverLocation{}
        drivers := []DriverInformation{}
        // Get listDriverLoc
        db.Find(&listDriverLoc)
        for _, data := range listDriverLoc {
            distance := helper.GetDistance(clientX, clientY, data.X, data.Y)
            if int(distance) <= threshold {
                driver := CreateDriverInformation(data, int(distance))
                drivers = append(drivers, driver)
            }
        }
        return drivers, nil
    } else {
        return nil, err
    }
}

func ChangeAvailable(id int, flag bool) (error) {
    if db, err := ConnectDatabase(); err == nil {
        defer db.Close()
        var driver DriverLocation
        db.Where("id = ?", id).First(&driver)
        // Check if empty
        driver.Available = &flag
        db.Save(&driver)
        return nil
    } else {
        return err
    }
}