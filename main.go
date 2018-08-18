package main

import (
    "app.goride/router"
    "github.com/gin-gonic/gin"
)

func main() {
    app := gin.Default()
    router.AddRoutesDriver(app)
    router.AddRoutesOrder(app)
    // Testing
    app.GET("/", func (c *gin.Context) {
        c.JSON(200, gin.H{"message": "hello"})
    })
    app.Run(":8080")
}