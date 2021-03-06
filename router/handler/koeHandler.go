package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/service"
	"github.com/kpkym/koe/utils/koe"
	"net/http"
)

var (
	pageSize = 12
)

func InitKoeHandler(group *gin.RouterGroup) {
	group.GET("/works/:category/*content", func(c *gin.Context) {
		pageRequest := dto.PageRequest{Size: pageSize}
		c.ShouldBindQuery(&pageRequest)

		works, count := service.NewService().WorkPage(pageRequest, c.Param("category"), c.Param("content")[1:])
		pagination := dto.Pagination{CurrentPage: pageRequest.Page, PageSize: pageRequest.Size, TotalCount: count}
		c.JSON(http.StatusOK, dto.SearchResponse{
			Pagination: pagination,
			Works:      works,
		})
	})

	group.GET("/search/*codes", func(c *gin.Context) {
		works := service.NewService().WorkCodes(koe.ListCode(c.Param("codes")))
		count := len(works)
		c.JSON(http.StatusOK, dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: count, TotalCount: int64(count)},
			Works:      works,
		})
	})

	group.GET("/work/:code", func(c *gin.Context) {
		work := service.NewService().WorkCodes([]string{c.Param("code")})[0]
		c.JSON(http.StatusOK, work)
	})

	group.PUT("/work/:code", func(c *gin.Context) {
		work := service.NewService().WorkCodes([]string{c.Param("code")})[0]
		c.JSON(http.StatusOK, work)
	})

	group.GET("/tracks/:code", func(c *gin.Context) {
		c.JSON(http.StatusOK, cache.NewMapCache[string, []*others.Node]().GetOrSet(c.Request.RequestURI, func() []*others.Node {
			return service.NewService().Track(c.Param("code"))
		}))
	})

	group.GET("/label/:category/", func(c *gin.Context) {
		c.JSON(http.StatusOK, service.NewService().Labels(c.Param("category")))
	})
}
