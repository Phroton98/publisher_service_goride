// This package is for defining Route for Order Controller
package router

import (
    "app.goride/app/controller/order"
    "app.goride/app/middleware"
    "github.com/gin-gonic/gin"
)

func AddRoutesOrder(route *gin.Engine) {
    // Set group
    group := route.Group("/api")
    group.Use(middleware.AuthRequired)
    {
        group.POST("/order", order.CreateOrder)
        group.GET("/order/:id", order.GetOrder)
        group.PUT("/order/:id", order.UpdateOrder)
        group.DELETE("/order/:id", order.CancelOrder)
        group.POST("/order/decline/:id", order.DeclineOrder)
    }
}