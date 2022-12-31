package utils

import "github.com/segmentio/ksuid"

type RandomUtils struct {
}

func NewRandomUtils() IRandomUtils {
	return &RandomUtils{}
}

func (ru *RandomUtils) GenerateUniqueID() string {
	return ksuid.New().String()
}
