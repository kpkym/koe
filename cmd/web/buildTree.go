package web

import (
	"github.com/google/uuid"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
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
	scanDir := global.GetServiceContext().Config.ScanDir

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

		if info.IsDir() {
			node.Type = utils.FolderType
			node.Children = make([]*others.Node, 0)
		} else {
			uuid := strings.Replace(uuid.NewString(), "-", "", -1)
			node.Type = getType(filepath.Ext(path))
			if codes := utils.ListCode(path); len(codes) > 0 {
				node.Code = codes[0]
			}
			node.UUID = uuid
			cache.NewMapCache[string]().Set(uuid, path)
		}
		parents[path] = node
		return nil
	}

	filepath.Walk(scanDir, walkFunc)
}

func BuildTree() []*others.Node {
	var result []*others.Node
	parents := make(map[string]*others.Node)
	scan(parents)
	loadNas(parents)

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = append(result, node)
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

	return result
}
