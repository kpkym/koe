package service

import (
	"fmt"
	"kpk-koe/utils"
	"testing"
)

func TestTrack(t *testing.T) {
	track := NewService().Track("294632")

	fmt.Println(utils.Marshal(track))

}
