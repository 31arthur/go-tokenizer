package server

import (
	"net/http"
	"server/internal/helper"
	"server/internal/utility"

	"github.com/gin-gonic/gin"
)

func (s *Server) helloWorld(c *gin.Context) {
	resp := helper.Response{
		Status: http.StatusOK,
		Data:   map[string]interface{}{"data": "Hello World to you too"},
	}
	helper.WriteJSONResponse(c.Writer, c.Request, resp)
}

func (s *Server) getUserDeets(c *gin.Context) {
	var user helper.User

	// Bind the JSON payload to the User struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := s.db.GetAccessTokenByID(user.ID)
	if err != nil {
		helper.ErrorResponse(
			"error in fetching access token @ "+err.Error(), c)
	}

	utility.SetAuthorizationCookie("pmt_auth_acct", accessToken, c)
}
