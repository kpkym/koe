package web

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	config := initConfig()

	fmt.Println(config)
	fmt.Println(config.Common)
}
