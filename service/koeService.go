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

func (s *service) Work(code string) domain.DlDomain {
	model := &domain.DlDomain{}

	koeDB := db.NewKoeDB[domain.DlDomain]()
	koeDB.GetData(model, code, func() domain.DlDomain {
		if c, err := colly.C(code); err == nil {
			marshal, _ := json.Marshal(c)
			return domain.DlDomain{Code: code, Data: string(marshal)}
		}
		return domain.DlDomain{Code: "0", Data: ""}
	})

	return *model
}

func (s *service) Track(code string) []others.Node {
	cacheHolder := pb.PBNode{
		Children: make([]*pb.PBNode, 0),
	}
	var resp []others.Node

	treeCacheKey := "tree"
	if cache.GetOrSet[*pb.PBNode](treeCacheKey, &cacheHolder, func() *pb.PBNode {
		logrus.Info("缓存为空 初始化目录树")

		copier.Copy(&cacheHolder.Children, utils.BuildTree())
		return &cacheHolder
	}) {
		logrus.Infof("缓存命中: %s", treeCacheKey)
	}

	copier.Copy(&resp, cacheHolder.GetChildren())

	trackCacheKey := fmt.Sprintf("track:%s", code)
	if cache.GetOrSet(trackCacheKey, &cacheHolder, func() *pb.PBNode {
		logrus.Infof("缓存为空 获取目录树: %s", code)

		tree := utils.GetTree(code, resp)
		copier.Copy(&cacheHolder.Children, tree)
		return &cacheHolder
	}) {
		logrus.Infof("缓存命中: %s", trackCacheKey)
	}

	copier.Copy(&resp, cacheHolder.GetChildren())
	return resp
}

func (s *service) GetFileFromUUID(code, uuid string) string {
	trackCacheKey := fmt.Sprintf("track:%s", code)
	cacheHolder := pb.PBNode{}
	if err := cache.Get(trackCacheKey, &cacheHolder); err != nil {
		logrus.Errorf("uuid获取, 缓存为空 获取目录树: %s", code)
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

func (s *service) GetLrcFromAudioUUID(code, uuid string) string {
	trackCacheKey := fmt.Sprintf("track:%s", code)
	cacheHolder := pb.PBNode{}
	if err := cache.Get(trackCacheKey, &cacheHolder); err != nil {
		logrus.Errorf("uuid获取, 缓存为空 获取目录树: %s", code)
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
