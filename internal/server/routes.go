package server

import (
	"net/http"
	"server/internal/middleware"
	"server/internal/tokenizer"
	"server/internal/utility"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	r.GET("/auth/:provider/callback", s.getAuthCallbackFunction)
	r.GET("/auth/:provider", s.getAuthComplete)
	r.GET("/logout/:provider", s.getLoggedOut)

	r.GET("/hello-world", middleware.AuthCheck(s.db), s.helloWorld)
	r.POST("/get-user-deets", s.getUserDeets)

	r.GET("/token-check", s.tokenHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) tokenHandler(c *gin.Context) {
	token, err := tokenizer.AccessTokenGenerator("hello")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to create token",
		})
		return
	}

	utility.SetAuthorizationCookie("Auth_Sample", token, c)
}
