package koe

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
)

func getType(ext string) string {
	switch strings.ToLower(ext) {
	case ".mp3", ".mp4", ".flac", ".wav", ".m4a":
		return "audio"
	case ".lrc", ".txt":
		return "text"
	case ".tif", ".jpg", ".jpeg", ".ico", ".tiff", ".gif", ".svg", ".webp", ".png", ".bmp":
		return "image"
	}
	return "other"
}

func scan(parents map[string]*others.Node) {
	scanDirs := global.GetServiceContext().Settings.ScanDirs

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error(err)
			return err
		}

		if info.Name() == ".DS_Store" {
			return nil
		}

		node := &others.Node{
			Title: info.Name(),
		}

		// uuid := strings.Replace(uuid.NewString(), "-", "", -1)
		uuid := utils.NextUUID()
		if info.IsDir() {
			node.Type = utils.FolderType
			node.Children = make([]*others.Node, 0)
		} else {
			node.Type = getType(filepath.Ext(path))
			if codes := ListCode(path); len(codes) > 0 {
				node.Code = codes[0]
			}
		}
		node.UUID = uuid
		cache.NewMapCache[uint32, string]().Set(uuid, path)
		parents[path] = node
		return nil
	}

	for _, e := range gjson.Parse(scanDirs.String()).Get("#.path").Array() {
		filepath.Walk(utils.IgnoreErr(homedir.Expand(e.String())), walkFunc)
	}
}

func BuildTree() {
	logrus.Info("初始化目录树")

	var trees []*others.Node
	parents := make(map[string]*others.Node)
	scan(parents)
	loadNas(parents)

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			trees = append(trees, node)
		} else {
			parent.Children = append(parent.Children, node)
			sort.Slice(parent.Children, func(i, j int) bool {
				if parent.Children[i].Type == utils.FolderType && parent.Children[j].Type == utils.FolderType {
					return strings.Compare(parent.Children[i].Title, parent.Children[j].Title) < 0
				} else if parent.Children[i].Type == utils.FolderType {
					return true
				} else if parent.Children[j].Type == utils.FolderType {
					return false
				}

				return strings.Compare(parent.Children[i].Title, parent.Children[j].Title) < 0
			})
		}
	}

	cache.NewMapCache[string, []*others.Node]().Set("trees", trees)
}

// Deprecated
func RemoveNotInTrees(tree []*others.Node, fn func([]string)) {
	utils.GoSafe(func() {
		table := global.GetServiceContext().DB.Table("work_domains").Session(&gorm.Session{})
		imageDir := filepath.Join(global.DataDir, "imgs")
		// 不在目录树的 需要删除db数据和图片
		var dbCodes []string
		table.Select("code").Scan(&dbCodes)
		fileTreeCodesSet := hashset.New(utils.Map[string, any](ListMyCode(tree), utils.Str2Any)...)
		dbCodesSet := hashset.New(utils.Map[string, any](dbCodes, utils.Str2Any)...)

		needCrawlCodesSet := fileTreeCodesSet.Difference(dbCodesSet).Values()
		if needRemoveCodesSet := dbCodesSet.Difference(fileTreeCodesSet).Values(); len(needRemoveCodesSet) > 0 {
			// 需要删除的
			table.Where("code IN ?", needRemoveCodesSet).Delete(&domain.WorkDomain{})
			for _, f := range utils.IgnoreErr(ioutil.ReadDir(imageDir)) {
				for _, rmCode := range needRemoveCodesSet {
					if strings.Contains(f.Name(), fmt.Sprint(rmCode)) {
						os.Remove(filepath.Join(imageDir, f.Name()))
					}
				}
			}
		}

		if needCrawlCodes := utils.Map[any, string](needCrawlCodesSet, utils.Any2Str); len(needCrawlCodes) > 0 {
			logrus.Info("爬虫抓取", needCrawlCodes)
			fn(needCrawlCodes)
		}
		logrus.Info("爬虫完成")
	}, "爬虫出错: ", string(debug.Stack()))
}

// Deprecated
func RemoveIncomplete() {
	imageDir := filepath.Join(global.DataDir, "imgs")

	// 检查db图片完整性 图片小于3个需要删除db数据重新抓取
	keys := make(map[string]int)
	var lt3Codes []string
	for _, entry := range utils.IgnoreErr(ioutil.ReadDir(imageDir)) {
		if cs := ListCode(entry.Name()); len(cs) != 0 {
			keys[cs[0]]++
		}
	}
	for k, v := range keys {
		if v != 3 {
			lt3Codes = append(lt3Codes, k)
		}
	}
	global.GetServiceContext().DB.Table("work_domains").Where("code IN ?", lt3Codes).Delete(&domain.WorkDomain{})
}
