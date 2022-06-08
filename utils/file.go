package utils

import (
	"github.com/google/uuid"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/others"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	folderType = "folder"
)

func BuildTree(fn func(string, string)) []*others.Node {
	var result []*others.Node

	scanDir := global.GetServiceContext().Config.ScanDir
	parents := make(map[string]*others.Node)

	getType := func(ext string) string {
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
			node.Type = folderType
			node.Children = make([]*others.Node, 0)
		} else {
			uuid := strings.Replace(uuid.NewString(), "-", "", -1)
			node.Type = getType(filepath.Ext(path))
			if codes := ListCode(path); len(codes) > 0 {
				node.Code = codes[0]
			}
			node.UUID = uuid
			fn(uuid, path)
		}
		parents[path] = node
		return nil
	}

	filepath.Walk(scanDir, walkFunc)

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = append(result, node)
		} else {
			parent.Children = append(parent.Children, node)

			sort.Slice(parent.Children, func(i, j int) bool {
				if parent.Children[i].Type == folderType && parent.Children[j].Type == folderType {
					return strings.Compare(parent.Children[i].Title, parent.Children[j].Title) < 0
				} else if parent.Children[i].Type == folderType {
					return true
				} else if parent.Children[j].Type == folderType {
					return false
				}

				return strings.Compare(parent.Children[i].Title, parent.Children[j].Title) < 0
			})
		}
	}

	return result
}

func GetTree(code string, tree []*others.Node) []*others.Node {
	nodes := make([]*others.Node, 0)
	getTreeHelper(code, &nodes, tree)
	return nodes
}

func getTreeHelper(code string, result *[]*others.Node, tree []*others.Node) {
	for _, e := range tree {
		if strings.Contains(e.Title, code) {
			*result = append(*result, e)
		} else {
			nodes := make([]*others.Node, 0)
			var c = &nodes
			for _, cc := range e.Children {
				*c = append(*c, cc)
			}
			getTreeHelper(code, result, *c)
		}
	}
}

func FlatTree(nodes []*others.Node) []*others.Node {
	fileNodes := make([]*others.Node, 0)

	for _, e := range nodes {
		if e.Type == folderType {
			for _, i := range FlatTree(e.Children) {
				fileNodes = append(fileNodes, i)
			}
		}
		fileNodes = append(fileNodes, e)
	}

	return fileNodes
}

func ScanToPo(scanPath string) []others.Po {
	pos := make([]others.Po, 0)

	filepath.Walk(IgnoreErr(homedir.Expand(scanPath)),
		func(path string, info os.FileInfo, err error) error {
			var filetype string

			if info.IsDir() {
				filetype = "D"
			} else {
				filetype = "F"
			}

			po := others.Po{
				Type: filetype,
				Ext:  strings.Replace(filepath.Ext(path), ".", "", 1),
				Path: path,
			}

			pos = append(pos, po)

			return nil
		})

	return pos
}

func GetFileBaseOnPwd(filepaths ...string) string {
	wd := IgnoreErr(os.Getwd())
	filepaths = append([]string{wd}, filepaths...)

	return filepath.Join(filepaths...)
}
