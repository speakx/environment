package rocksdbimp

import (
	"environment/mmapcache"

	"github.com/golang/protobuf/proto"
)

type Chunk interface {
	Append(pmsg proto.Message) bool
	Marshal() ([]byte, error)
	Unmarshal([]byte)
	NeedStore() bool
	GetChunkItem(i int) proto.Message
	Release()
}

type baseChunk struct {
	idx       int
	size      int
	mmapCache *mmapcache.MMapCache
}

func (b *baseChunk) initBaseChunk(size int) {
	b.size = size
	b.mmapCache = mmapcache.DefMMapCachePool.Alloc()
}

func (b *baseChunk) NeedStore() bool {
	if b.idx >= b.size {
		return true
	}
	return false
}

func (b *baseChunk) Release() {
	b.mmapCache.Release()
}
