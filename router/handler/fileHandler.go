package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/utils"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func InitFileHandler(group *gin.RouterGroup) {
	group.GET("/cover/:type/:id", func(c *gin.Context) {
		imgPath := filepath.Join(utils.GetFileBaseOnPwd("data, imgs"),
			filepath.Base(utils.GetImgUrl(c.Param("id"), c.Param("type"))))

		logrus.Infof("查找图片: %s", imgPath)
		c.File(imgPath)
	})
}
