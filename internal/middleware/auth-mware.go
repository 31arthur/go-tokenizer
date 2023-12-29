package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func validateAuth(c *gin.Context) {
	fmt.Println("Validating Auth")
}
