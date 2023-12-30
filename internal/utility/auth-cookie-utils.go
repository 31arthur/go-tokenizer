package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetAuthorizationCookie(key string, val string, c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	timeE, err := GetTokenExpiration(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error in getting token expiration date" + err.Error(),
		})
	}
	c.SetCookie(key, val, timeE.Second(), "/", "", false, true)
	// fmt.Println("Cookie Set")
	c.JSON(http.StatusOK, gin.H{
		"data": "Received user ID: " + val,
	})
}
