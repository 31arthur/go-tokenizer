package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUserToken(c *gin.Context, token string) {

	// Construct a JavaScript snippet to set the token in local storage
	js := `
    <script>
        // Set the token in local storage
        localStorage.setItem('accessToken', '` + token + `');
        console.log('Token set in local storage');
    </script>
    `

	// Set the Content-Type header to text/html
	c.Header("Content-Type", "text/html")

	// Send the JavaScript as part of the response
	c.String(http.StatusOK, js)
}
