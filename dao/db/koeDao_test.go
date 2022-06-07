package db

import (
	"encoding/json"
	"fmt"
	"github.com/kpkym/koe/dao"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/service"
	"testing"
)

func TestName(t *testing.T) {
}

func TestTrack(t *testing.T) {
	track := service.NewService().Track("360052")
	marshal, _ := json.Marshal(track)

	fmt.Println(string(marshal))
}

func TestCreateData(t *testing.T) {
	model := &domain.WorkDomain{}

	var newDB dao.Dao[domain.WorkDomain] = NewKoeDB[domain.WorkDomain]()

	err := newDB.GetData(model, "zzz", func() (domain.WorkDomain, error) {
		return domain.WorkDomain{Code: 123}, nil
	})

	fmt.Println(err)
	fmt.Println(model)

}
