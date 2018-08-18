// Package controller driver
package driver

import (
    "strings"
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
    // Variable
    var data driver.Location
    driverID := c.Param("id")
    // Check if request is right
    if err := c.ShouldBindJSON(&data); err == nil {
        // Create new data
        if err = driver.SetLocation(data, driverID); err == nil {
            c.JSON(http.StatusOK, gin.H{
                "status": strings.Join([]string{"user", driverID, "location updated"}, " "),
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
    }
}