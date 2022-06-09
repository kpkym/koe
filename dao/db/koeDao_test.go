package db

import (
	"encoding/json"
	"fmt"
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
}
