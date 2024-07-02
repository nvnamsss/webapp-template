package controllers

import (
	"webapp-template/src/utils"

	"github.com/gin-gonic/gin"
)

type Base struct {
}

// JSON responds a HTTP request with JSON data.
func (h *Base) JSON(c *gin.Context, data interface{}) {
	utils.JSON(c, data)
}

// HandleError handles error of HTTP request.
func (h *Base) HandleError(c *gin.Context, err error) {
	utils.HandleError(c, err)
}
