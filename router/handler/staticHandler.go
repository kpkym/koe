package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/global"
	"net/http"
)

func InitStaticHandler(group *gin.RouterGroup) {
	group.StaticFS("", http.Dir(global.GetServiceContext().Config.FlagConfig.ScanDir))
}
