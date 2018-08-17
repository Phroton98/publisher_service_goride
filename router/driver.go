// This package is for defining Route for Driver Controller
package router

import (
    "app.goride/app/controller/driver"
    "app.goride/app/middleware"
    "github.com/gin-gonic/gin"
)

func AddRoutesDriver(route *gin.Engine) {
    // Set group for driver
    group := route.Group("/api")
    group.Use(middleware.AuthRequired)
    {
        group.POST("/location/:id", driver.UpdateLocation)
        group.GET("/location/:id", driver.GetLocation)
    }
}