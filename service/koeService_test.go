package service

import (
	"fmt"
	"github.com/kpkym/koe/utils"
	"testing"
)

func TestTrack(t *testing.T) {
	track := NewService().Track("294632")

	fmt.Println(utils.Marshal(track))

}
