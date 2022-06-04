package service

import (
	"encoding/json"
	"github.com/kpkym/koe/colly"
	"github.com/kpkym/koe/dao/db"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
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

func (s *service) Track(id string) domain.TrackDomain {
	model := &domain.TrackDomain{}
	db.NewKoeDB[domain.TrackDomain]().GetData(model, id, func() domain.TrackDomain {
		fileSystemTree := &domain.FileSystemTreeDomain{}
		db.NewKoeDB[domain.FileSystemTreeDomain]().GetData(fileSystemTree, "treeAbsFalse", func() domain.FileSystemTreeDomain {
			tree := utils.BuildTree()
			marshal, _ := json.Marshal(tree)
			return domain.FileSystemTreeDomain{Code: "treeAbsFalse", Data: string(marshal)}
		})
		var tree []others.Node
		utils.Unmarshal(fileSystemTree.Data, &tree)

		getTree := utils.GetTree(id, tree)
		marshal, _ := json.Marshal(getTree)
		return domain.TrackDomain{Code: id, Data: string(marshal)}
	})
	return *model
}
