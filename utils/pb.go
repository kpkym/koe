package utils

import (
	"github.com/golang/protobuf/proto"
)

func PBUnmarshal[T proto.Message](data []byte, model T) {
	proto.Unmarshal(data, model)
}

func PBMarshal[T proto.Message](model T) []byte {
	return IgnoreErr(proto.Marshal(model))
}
