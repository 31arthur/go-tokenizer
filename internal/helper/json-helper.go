package helper

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status   int         `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	Redirect string      `json:"redirect,omitempty"`
	Error    error       `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, r *http.Request, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	if resp.Redirect != "" {
		http.Redirect(w, r, resp.Redirect, resp.Status)
		return
	}

	if err := json.NewEncoder(w).Encode(resp.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ErrorResponse(errStr string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errStr,
	})
}
