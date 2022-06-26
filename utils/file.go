package utils

import (
	"github.com/kpkym/koe/model/others"
	"github.com/mitchellh/go-homedir"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	FolderType = "folder"
)

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
			var c = nodes
			for _, cc := range e.Children {
				c = append(c, cc)
			}
			getTreeHelper(code, result, c)
		}
	}
}

func FlatTree(nodes []*others.Node) []*others.Node {
	fileNodes := make([]*others.Node, 0)

	for _, e := range nodes {
		if e.Type == FolderType {
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

	filepath.WalkDir(IgnoreErr(homedir.Expand(scanPath)),
		func(path string, info os.DirEntry, err error) error {
			var filetype string

			if info.IsDir() {
				filetype = "D"
			} else {
				filetype = "F"
			}

			po := others.Po{
				Type: filetype,
				// Ext:  strings.Replace(filepath.Ext(path), ".", "", 1),
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

func FilePath2Struct[V interface{}](filePath string) []V {
	list := make([]V, 0)
	value := gjson.ParseBytes(IgnoreErr(ioutil.ReadFile(IgnoreErr(homedir.Expand(filePath)))))

	value.ForEach(func(_, v gjson.Result) bool {
		var po V
		Unmarshal(v.Raw, &po)
		list = append(list, po)
		return true
	})

	return list
}
