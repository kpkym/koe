package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/service"
	"github.com/kpkym/koe/utils"
)

func InitKoeHandler(group *gin.RouterGroup) {
	group.GET("/works", func(c *gin.Context) {
		pageRequest := dto.PageRequest{Size: 5}
		c.ShouldBindQuery(&pageRequest)

		works, count := service.NewService().WorkPage(pageRequest)

		response := dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: pageRequest.Size, TotalCount: count},
			Works:      works,
		}

		c.JSON(200, response)
	})

	group.GET("/search/:code", func(c *gin.Context) {
		var works []domain.WorkDomain
		for _, code := range utils.ListCode(c.Param("code")) {
			works = append(works, service.NewService().WorkCode(code))
		}

		count := len(works)
		response := dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: count, TotalCount: int64(count)},
			Works:      works,
		}
		c.JSON(200, response)
	})

	group.GET("/work/:code", func(c *gin.Context) {
		work := service.NewService().WorkCode(c.Param("code"))
		c.JSON(200, work)
	})

	group.GET("/tracks/:code", func(c *gin.Context) {
		nodes := service.NewService().Track(c.Param("code"))
		c.JSON(200, nodes)
	})
}
