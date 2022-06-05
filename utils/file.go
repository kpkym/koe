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

func BuildTree() []others.Node {
	var result []*others.Node

	absRoot := global.GetServiceContext().Config.ScanDir
	parents := make(map[string]*others.Node)

	serve := global.GetServiceContext().Config.FlagConfig.Serve
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

		getType := func(ext string) string {
			if strings.Contains(ext, "mp3") ||
				strings.Contains(ext, "mp4") {
				return "audio"
			}
			return "file"
		}

		if info.IsDir() {
			node.Type = folderType
			node.Children = make([]*others.Node, 0)
		} else {
			serveFilePath := filepath.Join(serve, "static", path[len(absRoot):])
			node.Type = getType(filepath.Ext(serveFilePath))
			node.UUID = strings.Replace(uuid.NewString(), "-", "", -1)
			node.Duration = 1
			node.Path = path
		}
		parents[path] = node
		return nil
	}

	filepath.Walk(absRoot, walkFunc)

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

	var re []others.Node
	for _, r := range result {
		re = append(re, *r)
	}

	return re
}

func GetTree(code string, tree []others.Node) []others.Node {
	nodes := make([]others.Node, 0)
	var result = &nodes
	getTreeHelper(code, result, &tree)
	return *result
}

func getTreeHelper(code string, result *[]others.Node, tree *[]others.Node) {
	for _, e := range *tree {
		if strings.Contains(e.Title, code) {
			*result = append(*result, e)
		} else {
			nodes := make([]others.Node, 0)
			var c = &nodes
			for _, cc := range e.Children {
				*c = append(*c, *cc)
			}
			getTreeHelper(code, result, c)
		}
	}
}

func FlatTree(nodes []others.Node) []others.Node {
	fileNodes := make([]others.Node, 0)

	for _, e := range nodes {
		if e.Type == folderType {
			for _, i := range FlatTree(Map[*others.Node, others.Node](e.Children, func(item *others.Node) others.Node {
				return *item
			})) {
				fileNodes = append(fileNodes, i)
			}
		} else {
			fileNodes = append(fileNodes, e)
		}
	}

	return fileNodes
}

func ScanDir(scanPath string) []others.Po {
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

	filepaths = append([]string{filepath.Dir(wd)}, filepaths...)

	return filepath.Join(filepaths...)
}
