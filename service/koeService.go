package service

import (
	"fmt"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/dao/db"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"path/filepath"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) WorkPage(pageRequest dto.PageRequest, category, content string) ([]domain.WorkDomain, int64) {
	var resp []domain.WorkDomain

	fn := func(db *gorm.DB) {
		switch category {
		case "circle":
			db.Where(fmt.Sprintf("%s = ?", category), content)
		case "vas", "tags":
			db.Joins(fmt.Sprintf("join json_each(%s) j", category)).Where("j.value = ?", content)
		}
	}

	return resp, db.NewKoeDB[domain.WorkDomain](domain.WorkDomain{}).Page(&resp, pageRequest, fn)
}

func (s *service) WorkCodes(codes []string) []domain.WorkDomain {
	var resp []domain.WorkDomain
	db.NewKoeDB[[]domain.WorkDomain](domain.WorkDomain{}).ListByCode(&resp, codes)

	return resp
}

func (s *service) Labels(category string) *[]dto.LabelResponse {

	return cache.NewMapCache[*[]dto.LabelResponse]().GetOrSet(fmt.Sprintf("category:%s", category), func() *[]dto.LabelResponse {
		var results = new([]dto.LabelResponse)
		var fn func(g *gorm.DB)
		switch category {
		case "vas", "tags":
			fn = func(g *gorm.DB) {
				g.Select("j.value name, COUNT(1) count").
					Joins(fmt.Sprintf("join json_each(%s) j", category)).
					Group("j.value").Order("count DESC")
			}
		default:
			fn = func(g *gorm.DB) {
				g.Select(fmt.Sprintf("%s name, COUNT(1) count", category)).Group(category).Order("count DESC")
			}
		}
		db.NewKoeDB[[]domain.WorkDomain](domain.WorkDomain{}).RunRow(results, fn)
		return results
	})
}

func (s *service) Track(code string) []*others.Node {
	tree, _ := cache.NewMapCache[[]*others.Node]().Get("tree")

	return cache.NewMapCache[[]*others.Node]().GetOrSet(fmt.Sprintf("track:%s", code), func() []*others.Node {
		logrus.Infof("缓存为空 获取目录树: %s", code)
		getTree := utils.GetTree(code, tree)
		return getTree
	})
}

func (s *service) GetFileFromUUID(uuid string) string {
	value, _ := cache.NewMapCache[string]().Get(uuid)
	return value
}

func (s *service) GetLrcFromAudioUUID(code, uuid string) string {
	trackCacheKey := fmt.Sprintf("track:%s", code)
	audioFilePath, _ := cache.NewMapCache[string]().Get(uuid)
	nodes, _ := cache.NewMapCache[[]*others.Node]().Get(trackCacheKey)

	lrcPath, _ := utils.GetLrcPath(filepath.Base(audioFilePath), nodes, func(uuid string) string {
		lrcFilepath, _ := cache.NewMapCache[string]().Get(uuid)
		return lrcFilepath
	})
	logrus.Infof("查找文件: %s的lrc文件. 结果为为: %s", audioFilePath, lrcPath)
	return lrcPath
}
