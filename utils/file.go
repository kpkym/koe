package utils

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"kpk-koe/global"
	"kpk-koe/model/others"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func BuildTree() []others.Node {
	var result []*others.Node
	absRoot := global.ScanDir
	parents := make(map[string]*others.Node)
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
			node.Type = "folder"
			node.Children = make([]*others.Node, 0)
		} else {
			servePath := fmt.Sprintf("http://127.0.0.1:8081/static%s", path[len(global.ScanDir):])
			node.Type = getType(filepath.Ext(servePath))
			node.Hash = servePath
			node.WorkTitle = info.Name()
			node.MediaStreamUrl = servePath
			node.MediaDownloadUrl = servePath
			node.LrcUrl = strings.Replace(servePath, ".mp3", ".lrc", 1)
			node.Duration = 1

			if rjs := ListRJ(servePath); len(rjs) > 0 {
				node.ImgUrl = fmt.Sprintf("", rjs[0])
			}
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
				if parent.Children[i].Type == "folder" && parent.Children[j].Type == "folder" {
					return strings.Compare(parent.Children[i].Title, parent.Children[j].Title) < 0
				} else if parent.Children[i].Type == "folder" {
					return true
				} else if parent.Children[j].Type == "folder" {
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
