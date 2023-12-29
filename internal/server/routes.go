package server

import (
	"net/http"

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

	r.GET("/hello-world", s.helloWorld)
	r.POST("/get-user-deets", s.getUserDeets)

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
