package db

import (
	"encoding/json"
	"fmt"
	"kpk-koe/dao"
	"kpk-koe/model/domain"
	"kpk-koe/service"
	"testing"
)

func TestName(t *testing.T) {
}

func TestTrack(t *testing.T) {
	track := service.NewService().Track("RJ360052")
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
