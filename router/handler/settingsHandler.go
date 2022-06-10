package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitSettingsHandler(group *gin.RouterGroup) {
	group.GET("/scan", func(c *gin.Context) {

		c.JSON(http.StatusOK, "ok")
	})
}
