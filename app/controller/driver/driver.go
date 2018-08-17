// Package controller driver
package driver

import (
    "strings"
    "strconv"
    "net/http"
    "app.goride/model/driver"
    "github.com/gin-gonic/gin"
)

func GetLocation(c *gin.Context) {
    driverID := c.Param("id")
    if json, err := driver.GetLocation(driverID); err == nil {
        c.JSON(http.StatusOK, json)
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}

func UpdateLocation(c *gin.Context) {
    driverID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
    var data driver.Location
    if err = c.ShouldBindJSON(&data); err == nil {
        // Create new data
        newData := driver.CreateDriverLocation(driverID, data)
        if err := driver.SetLocation(&newData); err == nil {
            c.JSON(http.StatusOK, gin.H{
                "status": strings.Join([]string{"user", strconv.Itoa(newData.ID), "location updated"}, " "),
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}