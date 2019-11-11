package mmapcache

import (
	"environment/byteio"
)

const (
	mmapdataEnable  uint8 = 0
	mmapdataDisable uint8 = 1
)

type mmapdata struct {
	size []byte
	flag []byte
	data []byte
}

func (m *mmapdata) recallMMapData() {
	byteio.Uint8ToBytes(mmapdataDisable, m.flag)
}
