package middleware

import "github.com/gin-gonic/gin"
import "fmt"

func AuthRequired(c *gin.Context) {
	fmt.Println("Hello")
	// To Do Check Authentication
	c.Next()
}