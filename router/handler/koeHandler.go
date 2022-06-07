package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/service"
	"github.com/kpkym/koe/utils"
)

func InitKoeHandler(group *gin.RouterGroup) {
	group.GET("/works", func(c *gin.Context) {
		pageRequest := dto.PageRequest{Size: 5}
		c.ShouldBindQuery(&pageRequest)
		works, count := service.NewService().WorkPage(pageRequest)

		c.JSON(200, dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: pageRequest.Size, TotalCount: count},
			Works:      works,
		})
	})

	group.GET("/search/:codes", func(c *gin.Context) {
		works := service.NewService().WorkCodes(utils.ListCode(c.Param("codes")))
		count := len(works)

		c.JSON(200, dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: count, TotalCount: int64(count)},
			Works:      works,
		})
	})

	group.GET("/work/:code", func(c *gin.Context) {
		work := service.NewService().WorkCodes([]string{c.Param("code")})[0]
		c.JSON(200, work)
	})

	group.GET("/tracks/:code", func(c *gin.Context) {
		nodes := service.NewService().Track(c.Param("code"))
		c.JSON(200, nodes)
	})
}
