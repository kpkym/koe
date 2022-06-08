package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/dao/db"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/model/pb"
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

func (s *service) Labels(category string) []dto.LabelResponse {
	var results []dto.LabelResponse

	cache.GetOrSetJSON[[]dto.LabelResponse](fmt.Sprintf("category:%s", category), &results, func() []dto.LabelResponse {
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
		db.NewKoeDB[[]domain.WorkDomain](domain.WorkDomain{}).RunRow(&results, fn)
		return results
	})

	return results
}

func (s *service) Track(code string) []*others.Node {
	var resp = cache.GetJSON[[]*others.Node]("tree")

	trackCacheKey := fmt.Sprintf("track:%s", code)
	if cache.GetOrSetJSON[[]*others.Node](trackCacheKey, &resp, func() []*others.Node {
		logrus.Infof("缓存为空 获取目录树: %s", code)
		return utils.GetTree(code, resp)
	}) {
		logrus.Infof("缓存命中: %s", trackCacheKey)
	}

	return resp
}

func (s *service) GetFileFromUUID(uuid string) string {
	return cache.GetJSON[string](uuid)
}

func (s *service) GetLrcFromAudioUUID(code, uuid string) string {
	trackCacheKey := fmt.Sprintf("track:%s", code)
	cacheHolder := pb.PBNode{}
	if err := cache.Get(trackCacheKey, &cacheHolder); err != nil {
		logrus.Errorf("uuid获取, 缓存为空 获取目录树: %s", code)
		return ""
	}

	var resp []*others.Node
	var filePath string
	copier.Copy(&resp, cacheHolder.GetChildren())
	for _, e := range utils.FlatTree(resp) {
		if e.UUID == uuid {
			filePath = cache.GetJSON[string](uuid)
			break
		}
	}

	if filePath == "" {
		logrus.Errorf("uuid获取, 返回为空")
		return ""
	}

	lrcPath, _ := utils.GetLrcPath(filepath.Base(filePath), resp, func(uuid string) string {
		return cache.GetJSON[string](uuid)
	})
	logrus.Infof("查找文件: %s的lrc文件. 结果为为: %s", filePath, lrcPath)
	return lrcPath
}
