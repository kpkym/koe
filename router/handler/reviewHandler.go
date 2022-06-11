package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"net/http"
)

func InitReviewHandler(group *gin.RouterGroup) {
	group.PUT("/review", func(c *gin.Context) {
		workDomain := new(domain.WorkDomain)
		c.ShouldBindJSON(workDomain)

		global.GetServiceContext().DB.Model(workDomain).Update("user_rating", workDomain.UserRating)
		c.JSON(http.StatusOK, "ok")
	})

}
