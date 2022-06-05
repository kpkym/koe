package service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/kpkym/koe/colly"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/dao/db"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/model/pb"
	"github.com/kpkym/koe/utils"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Work(id string) domain.DlDomain {
	model := &domain.DlDomain{}

	koeDB := db.NewKoeDB[domain.DlDomain]()
	koeDB.GetData(model, id, func() domain.DlDomain {
		if c, err := colly.C(id); err == nil {
			marshal, _ := json.Marshal(c)
			return domain.DlDomain{Code: id, Data: string(marshal)}
		}
		return domain.DlDomain{Code: "0", Data: ""}
	})

	return *model
}

func (s *service) Track(id string) []others.Node {
	cacheHolder := pb.PBNode{}
	var resp []others.Node

	treeCacheKey := "tree"
	if err := cache.Get(treeCacheKey, &cacheHolder); err != nil {
		logrus.Info("缓存为空 初始化目录树")

		var treePB []*pb.PBNode
		copier.Copy(&treePB, utils.BuildTree())
		cacheHolder.Children = treePB
		cache.Set(treeCacheKey, &cacheHolder)
	} else {
		logrus.Infof("缓存命中: %s", treeCacheKey)
	}

	copier.Copy(&resp, cacheHolder.GetChildren())

	trackCacheKey := fmt.Sprintf("track:%s", id)
	if err := cache.Get(trackCacheKey, &cacheHolder); err != nil {
		logrus.Infof("缓存为空 获取目录树: %s", id)

		var treePB []*pb.PBNode
		copier.Copy(&treePB, utils.GetTree(id, resp))
		cacheHolder.Children = treePB
		cache.Set(trackCacheKey, &cacheHolder)
	} else {
		logrus.Infof("缓存命中: %s", trackCacheKey)
	}

	copier.Copy(&resp, cacheHolder.GetChildren())
	return resp
}

func (s *service) GetFileFromUUID(id, uuid string) string {
	trackCacheKey := fmt.Sprintf("track:%s", id)
	cacheHolder := pb.PBNode{}
	if err := cache.Get(trackCacheKey, &cacheHolder); err != nil {
		logrus.Errorf("uuid获取, 缓存为空 获取目录树: %s", id)
		return ""
	}

	var resp []others.Node
	copier.Copy(&resp, cacheHolder.GetChildren())
	for _, e := range utils.FlatTree(resp) {
		if e.UUID == uuid {
			return e.Path
		}
	}

	logrus.Errorf("uuid获取, 返回为空")
	return ""
}

func (s *service) GetLrcFromAudioUUID(id, uuid string) string {
	trackCacheKey := fmt.Sprintf("track:%s", id)
	cacheHolder := pb.PBNode{}
	if err := cache.Get(trackCacheKey, &cacheHolder); err != nil {
		logrus.Errorf("uuid获取, 缓存为空 获取目录树: %s", id)
		return ""
	}

	var resp []others.Node
	var filePath string
	copier.Copy(&resp, cacheHolder.GetChildren())
	for _, e := range utils.FlatTree(resp) {
		if e.UUID == uuid {
			filePath = e.Path
			break
		}
	}

	if filePath == "" {
		logrus.Errorf("uuid获取, 返回为空")
		return ""
	}

	lrcPath, _ := utils.GetLrcPath(filepath.Base(filePath), resp)
	logrus.Infof("查找文件: %s的lrc文件. 结果为为: %s", filePath, lrcPath)
	return lrcPath
}
