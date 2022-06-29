package service

import (
	"fmt"
	dbf "github.com/kpkym/koe/cmd/web/config/db"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/dao/db"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"github.com/kpkym/koe/utils/koe"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
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
			db.Where(datatypes.JSONQuery(category).HasKey(content))
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

	return cache.NewMapCache[string, *[]dto.LabelResponse]().GetOrSet(fmt.Sprintf("category:%s", category), func() *[]dto.LabelResponse {
		var results = new([]dto.LabelResponse)
		var fn func(g *gorm.DB)
		switch category {
		case "vas", "tags":
			fn = func(g *gorm.DB) {
				dbf.SelectGroupBy(g, category)
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
	tree, _ := cache.NewMapCache[string, []*others.Node]().Get("trees")

	return cache.NewMapCache[string, []*others.Node]().GetOrSet(fmt.Sprintf("track:%s", code), func() []*others.Node {
		logrus.Infof("缓存为空 获取目录树: %s", code)
		getTree := utils.GetTree(code, tree)
		return getTree
	})
}

func (s *service) GetFileFromUUID(uuid uint32) string {
	value, _ := cache.NewMapCache[uint32, string]().Get(uuid)
	return value
}

func (s *service) GetLrcFromAudioUUID(code string, uuid uint32) uint32 {
	trackCacheKey := fmt.Sprintf("track:%s", code)
	audioFilePath, _ := cache.NewMapCache[uint32, string]().Get(uuid)
	nodes, _ := cache.NewMapCache[string, []*others.Node]().Get(trackCacheKey)

	fn := func(uuid uint32) string {
		lrcFilepath, _ := cache.NewMapCache[uint32, string]().Get(uuid)
		return lrcFilepath
	}
	lrcUUID, _ := koe.GetLrcUUID(filepath.Base(audioFilePath), nodes, fn)
	logrus.Infof("查找文件: %s的lrc文件, 结果UUID为: %v, 路径为: %v", audioFilePath, lrcUUID, fn(lrcUUID))
	return lrcUUID
}

func (s *service) GetLrcsFromAudioUUID(code string) []*others.Node {
	trackCacheKey := fmt.Sprintf("track:%s", code)
	nodes, _ := cache.NewMapCache[string, []*others.Node]().Get(trackCacheKey)

	return utils.Filter[*others.Node](utils.FlatTree(nodes), func(item *others.Node) bool {
		return filepath.Ext(item.Title) == ".lrc"
	})
}
