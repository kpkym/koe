package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/service"
	"net/http"
	"os/exec"
)

func InitSettingsHandler(group *gin.RouterGroup) {
	group.GET("/scan", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	group.GET("/open/:uuid", func(c *gin.Context) {
		filePath := service.NewService().GetFileFromUUID(c.Param("uuid"))
		exec.Command("open", fmt.Sprintf("%s", filePath)).Run()
		c.JSON(http.StatusOK, "ok")
	})

	settings := group.Group("/settings")
	settings.
		GET("/", func(c *gin.Context) {
			settings := &domain.Settings{ID: 1}
			global.GetServiceContext().DB.First(settings)
			c.JSON(http.StatusOK, settings)
		}).
		PUT("/", func(c *gin.Context) {
			settings := &domain.Settings{ID: 1}
			c.ShouldBindJSON(settings)
			global.GetServiceContext().DB.Save(settings)
			c.JSON(http.StatusOK, "ok")
		})
}
