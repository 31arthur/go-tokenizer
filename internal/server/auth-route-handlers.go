package server

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/helper"
	"server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func (s *Server) getAuthCallbackFunction(c *gin.Context) {

	provider := c.Param("provider")

	r := c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(c.Writer, r)
	if err != nil {
		fmt.Fprintln(c.Writer, r)
		resp := helper.Response{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		helper.WriteJSONResponse(c.Writer, r, resp)
		return
	}
	fmt.Println(user)

	userAuth := models.MapUserToUserAuth(user)
	if err = s.db.InsertNewUser(&userAuth); err != nil {
		fmt.Fprintln(c.Writer, r)
		resp := helper.Response{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		helper.WriteJSONResponse(c.Writer, r, resp)
		return
	}

	resp := helper.Response{
		Status:   http.StatusFound,
		Redirect: "http://localhost:5173/landing?user=" + user.UserID,
	}
	helper.WriteJSONResponse(c.Writer, r, resp)
}

func (s *Server) getAuthComplete(c *gin.Context) {

	provider := c.Param("provider")
	w := c.Writer
	r := c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) getLoggedOut(c *gin.Context) {
	provider := c.Param("provider")
	w := c.Writer
	r := c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
