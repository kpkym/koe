package koe

import (
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"path/filepath"
)

func loadNas(parents map[string]*others.Node, nasPrefix string) {
	for _, e := range GetNasJson() {
		name := filepath.Base(e.Path)
		path := filepath.Join(nasPrefix, e.Path)
		if name == ".DS_Store" {
			continue
		}

		node := &others.Node{Title: name}
		// uuid := strings.Replace(uuid.NewString(), "-", "", -1)
		uuid := utils.NextUUID()

		if e.Type == "D" {
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
	}
}
