package main

import (
    "app.goride/router"
    "github.com/gin-gonic/gin"
)

func main() {
    app := gin.Default()
    router.AddRoutesDriver(app)
    router.AddRoutesOrder(app)
    app.Run(":8080")
}