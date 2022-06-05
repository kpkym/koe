package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/service"
	"github.com/kpkym/koe/utils"
)

func InitKoeHandler(group *gin.RouterGroup) {
	group.GET("/works", func(c *gin.Context) {
		works := make([]dto.Work, 0)
		for _, code := range utils.ListMyCode() {
			item := dto.Work{}
			utils.Unmarshal(service.NewService().Work(code).Data, &item)
			works = append(works, item)
		}

		count := len(works)
		response := dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: count, TotalCount: count},
			Works:      works,
		}

		c.JSON(200, response)
	})

	group.GET("/search/:id", func(c *gin.Context) {
		works := dto.Work{}
		utils.Unmarshal(service.NewService().Work(c.Param("id")).Data, &works)

		response := dto.SearchResponse{
			Pagination: dto.Pagination{CurrentPage: 1, PageSize: 1, TotalCount: 1},
			Works:      []dto.Work{works},
		}
		c.JSON(200, response)
	})

	group.GET("/work/:id", func(c *gin.Context) {
		work := dto.Work{}
		utils.Unmarshal(service.NewService().Work(c.Param("id")).Data, &work)
		c.JSON(200, work)
	})

	group.GET("/tracks/:id", func(c *gin.Context) {
		nodes := service.NewService().Track(c.Param("id"))
		c.JSON(200, nodes)
	})
}
