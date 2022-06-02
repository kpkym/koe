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
	model := &domain.TrackDomain{}

	var newDB dao.Dao[domain.TrackDomain] = NewKoeDB[domain.TrackDomain]()

	err := newDB.GetData(model, "zzz", func() domain.TrackDomain {
		return domain.TrackDomain{Code: "zzz", Data: "mmm"}
	})

	fmt.Println(err)
	fmt.Println(model)

}
