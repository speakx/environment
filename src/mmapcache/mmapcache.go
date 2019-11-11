package mmapcache

import (
	"errors"
	"os"

	"github.com/edsrzf/mmap-go"

	"github.com/speakx/environment/byteio"
)

const (
	// FlushModeTimeInvalidate 通过激活时间来决定是否需要将有改变的数据Flush到DB
	// 例如，time.Second * 60 * 5 内，没有访问就Flush & 释放
	FlushModeTimeInvalidate = iota
	// FlushModeRecordCount 通过记录数量来决定是否Flush到DB
	// 例如，新建记录满50条后刷新到DB
	FlushModeRecordCount
)

const (
	mmapHeadSize = 28
)

// MMapCache 基于mmap模式的文件写对象
// | --------------------------------------- head --------------------------------------------------------- | ----------------- content ---------------|
// | - mmap.totalsize - | -------------------------- mmap.setting ------------------------------------------| ------ mmap.data ------ |-- mmap.data -- |
// | 4byte:content.size | 8byte:id | 8byte:next-id | 3byte:version | 1byte:flushmode | 4byte:flushCondition | 4byte:size | 1byte:flag | .............. |
type MMapCache struct {
	path             string
	f                *os.File
	buf              []byte // mmap后的文件原始内存
	writeContent     []byte // content部分的内存对象
	writeUint32Cache []byte // 内存写缓存，保证一次copy写内存，防止按字节写出错
	readPos          int
	writePos         int
	flushMode        uint8
	flushCondition   uint32
	activeTimeStamp  int64
	nextMMapCache    *MMapCache
	mmapdataIdx      map[string]*mmapdata
}

func newMMapCache(filePath string) (*MMapCache, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if nil != err {
		return nil, err
	}
	buf, err := mmap.MapRegion(f, -1, mmap.RDWR, 0, 0)
	if nil != err {
		return nil, err
	}

	mmcache := &MMapCache{
		path:             filePath,
		f:                f,
		buf:              buf,
		writeContent:     buf[mmapHeadSize:],
		writeUint32Cache: make([]byte, 4),
	}

	mmcache.init()
	return mmcache, nil
}

// GetID 获取当前mmap对象的ID
func (m *MMapCache) GetID() uint64 {
	return byteio.BytesToUint64(m.getIDBuf())
}

// GetNextID 当多片mmap合并到一起使用时获取到下一片的mmap对象
func (m *MMapCache) GetNextID() uint64 {
	return byteio.BytesToUint64(m.getNextIDBuf())
}

// MergeMMapCache 将另一片mmap与当前mmap合并到一起产生一个新的大块mmap对象
func (m *MMapCache) MergeMMapCache(mmapCache *MMapCache) *MMapCache {
	mmapCache.nextMMapCache = m
	byteio.Uint64ToBytes(m.GetID(), mmapCache.getNextIDBuf())
	return mmapCache
}

// Release 释放，将此mmap文件丢到pool中，由pool的策略决定释放真正释放
func (m *MMapCache) Release() {
	mmapCache := m
	for {
		next := mmapCache.nextMMapCache
		DefPoolMMapCache.Collect(mmapCache)

		if nil == next {
			break
		}
	}
}

// Write 写入一片内存对象，当key不为空时，意味着需要建立一个映射
// 可以通过key，随时把写入的这片内存作废
func (m *MMapCache) Write(p []byte, key []byte) (int, error) {
	if len(p) > m.GetFreeDataLen() {
		return 0, errors.New("mmap cache buf over follow")
	}

	copy(m.writeContent, p)
	m.setWritePos(len(p) + m.writePos)

	if len(key) > 0 {
		m.mmapdataIdx[string(key)] = &mmapdata{
			size: m.writeContent[m.writePos-len(p):],
			flag: m.writeContent[m.writePos-len(p)+4:],
			data: m.writeContent[m.writePos-len(p)+4+1:],
		}
	}

	return len(p), nil
}

// RecallWrite 通过key，把之前写入的缓存废弃掉（但这些缓存是不复用的）
func (m *MMapCache) RecallWrite(key string) bool {
	mmapdata, ok := m.mmapdataIdx[key]
	if false == ok || nil == mmapdata {
		return false
	}
	mmapdata.recallMMapData()
	return true
}

// GetFreeDataLen 获取剩下的可以write的数据空间
func (m *MMapCache) GetFreeDataLen() int {
	return len(m.writeContent) - m.writePos
}

func (m *MMapCache) name(id uint64) {
	byteio.Uint64ToBytes(id, m.getIDBuf())
}

func (m *MMapCache) close(remove bool) {
	m.f.Close()
	if remove {
		os.Remove(m.path)
	}
}

func (m *MMapCache) getIDBuf() []byte {
	return m.buf[4:]
}

func (m *MMapCache) getNextIDBuf() []byte {
	return m.buf[12:]
}

func (m *MMapCache) setWritePos(n int) {
	byteio.SafeUint32ToBytes(uint32(n), m.buf, m.writeUint32Cache)
	m.writePos = n
}

func (m *MMapCache) getWritePos() int {
	return int(byteio.BytesToUint32(m.buf))
}

func (m *MMapCache) init() {
	m.readPos = 0
	m.setWritePos(0)

	m.nextMMapCache = nil
	m.mmapdataIdx = make(map[string]*mmapdata)
}

func (m *MMapCache) recycle(template []byte) {
	m.init()
}
