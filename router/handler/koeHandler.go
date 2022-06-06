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
		works := make([]domain.WorkDomain, 0)
		for _, code := range utils.ListMyCode() {
			works = append(works, service.NewService().Work(code))
		}

		count := len(works)
		response := dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: count, TotalCount: count},
			Works:      works,
		}

		c.JSON(200, response)
	})

	group.GET("/search/:code", func(c *gin.Context) {
		works := service.NewService().Work(c.Param("code"))

		response := dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: 1, TotalCount: 1},
			Works:      []domain.WorkDomain{works},
		}
		c.JSON(200, response)
	})

	group.GET("/work/:code", func(c *gin.Context) {
		work := service.NewService().Work(c.Param("code"))
		c.JSON(200, work)
	})

	group.GET("/tracks/:code", func(c *gin.Context) {
		nodes := service.NewService().Track(c.Param("code"))
		c.JSON(200, nodes)
	})
}
