package utils

import (
	"encoding/json"
)

func Unmarshal[V any](d string, v V) {
	json.Unmarshal([]byte(d), v)
}

func Marshal[V any](v V) string {
	return string(IgnoreErr(json.Marshal(v)))
}
