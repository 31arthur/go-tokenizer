package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetAuthorizationCookie(key string, val string, c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(key, val, 3600*24, "", "", false, true)
	// fmt.Println("Cookie Set")
	c.JSON(http.StatusOK, gin.H{})
}
