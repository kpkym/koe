package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/service"
	"github.com/kpkym/koe/utils"
	"github.com/kpkym/koe/utils/koe"
	"github.com/sirupsen/logrus"
	"net/http"
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
		fileUUID := koeService.GetLrcFromAudioUUID(c.Param("code"), utils.Str2Num[uint32](c.Param("uuid")))
		c.String(http.StatusOK, utils.Any2Str(fileUUID))
	})

	group.GET("/lrcs/:code", func(c *gin.Context) {
		koeService := service.NewService()
		nodes := koeService.GetLrcsFromAudioUUID(c.Param("code"))

		c.JSON(http.StatusOK, utils.Map[*others.Node, gin.H](nodes, func(node *others.Node) gin.H {
			return gin.H{"uuid": node.UUID, "name": node.Title}
		}))
	})

	group.GET("/:uuid", func(c *gin.Context) {
		filePath := service.NewService().GetFileFromUUID(utils.Str2Num[uint32](c.Param("uuid")))
		logrus.Infof("查找文件: %s", filePath)
		c.File(filePath)
	})
}
