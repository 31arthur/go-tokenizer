package server

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/helper"
	"server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
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
	// fmt.Println(user)
	uid := s.ExtendAuthCallback(user, c)

	// fmt.Println(uid)

	resp := helper.Response{
		Status:   http.StatusFound,
		Redirect: "http://localhost:5173/landing?user=" + uid,
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

//========================================================================

func (s *Server) ExtendAuthCallback(user goth.User, c *gin.Context) string {
	flag, id, err := s.db.FreshCheck(user.UserID, user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		// Handle the database error
		helper.ErrorResponse("Error in fresh Check @ || "+err.Error(), c)
		return ""
	}

	userAuthDetails, err := models.MapUserToUserAuth(user, id)
	if err != nil {
		helper.ErrorResponse("Error mapping User to User Auth @ "+err.Error(), c)
		return ""
	}

	userDetails := models.MapUserToUserDetails(user, userAuthDetails.ID)

	if !flag {
		if err := s.db.InsertNewUser(&userAuthDetails); err != nil {
			helper.ErrorResponse("Error inserting new User Auth Table @ "+err.Error(), c)
			return ""
		}

		if err := s.db.InsertNewUser(&userDetails); err != nil {
			helper.ErrorResponse("Error inserting new User Details Table @ "+err.Error(), c)
			return ""
		}
		return userAuthDetails.ID
	}

	if err := s.db.ReplaceUserAuthByID(&userAuthDetails); err != nil {
		helper.ErrorResponse("Error updating User Auth Table @ "+err.Error(), c)
		return ""
	}

	return userAuthDetails.ID
}
