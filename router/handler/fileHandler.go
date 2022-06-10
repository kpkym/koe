package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/service"
	"github.com/kpkym/koe/utils/koe"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func InitFileHandler(group *gin.RouterGroup) {
	group.GET("/cover/:code/:type", func(c *gin.Context) {
		imgPath := filepath.Join(global.DataDir, "imgs",
			filepath.Base(koe.GetImgUrl(c.Param("code"), c.Param("type"))))
		logrus.Infof("查找图片: %s", imgPath)

		c.File(imgPath)
	})

	group.GET("/lrc/:code/:uuid", func(c *gin.Context) {
		koeService := service.NewService()
		filePath := koeService.GetLrcFromAudioUUID(c.Param("code"), c.Param("uuid"))

		c.File(filePath)
	})

	group.GET("/:uuid", func(c *gin.Context) {
		filePath := service.NewService().GetFileFromUUID(c.Param("uuid"))
		logrus.Infof("查找文件: %s", filePath)

		c.File(filePath)
	})
}
