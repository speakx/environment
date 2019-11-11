package rocksdbimp

import (
	envrockspb "environment/rocksdbimp/proto"

	"github.com/golang/protobuf/proto"
)

type ChunkMessage struct {
	baseChunk
	chunkKey []byte
}

func NewChunkMessage(size int) *ChunkMessage {
	c := &ChunkMessage{}
	c.baseChunk.initBaseChunk(size)
	return c
}

func (c *ChunkMessage) Append(pmsg proto.Message) bool {
	buf, _ := proto.Marshal(pmsg)
	c.mmapCache.Write(buf)
	c.chunkKey = pmsg.(*envrockspb.Message).Key
	return true
}

func (c *ChunkMessage) GetChunkKey() []byte {
	return c.chunkKey
}

func (c *ChunkMessage) Marshal() ([]byte, error) {
	// return proto.Marshal(c)
	return c.mmapCache.GetWriteData(), nil
}

func (c *ChunkMessage) Unmarshal([]byte) {
}
