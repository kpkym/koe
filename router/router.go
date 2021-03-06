package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/router/handler"
	"net/http"
)

func GetGinServe() *gin.Engine {
	var engin = gin.Default()

	// engin.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	//
	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(204)
	// 		return
	// 	}
	// 	c.Next()
	// })

	api := engin.Group("api")

	handler.InitKoeHandler(api)
	handler.InitSettingsHandler(api)
	handler.InitReviewHandler(api)
	handler.InitFileHandler(api)

	api.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	return engin
}
