package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TheTest(c *gin.Context) {
	fmt.Println("this is a middleware!")
}
