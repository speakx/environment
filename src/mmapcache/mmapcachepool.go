package mmapcache

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// DefPoolMMapCache Init初始化完成后，得到的内存池对象
var DefPoolMMapCache *PoolMMapCache

const mmapInitByte byte = 0

// PoolMMapCache 内存池对象
type PoolMMapCache struct {
	dir            string
	template       []byte
	recycleDur     time.Duration
	allocator      chan *MMapCache
	collector      chan *MMapCache
	errorfuc       func(error)
	loadFlag       bool
	allocCounter   uint64
	collectCounter uint64
	releaseCounter uint64
}

// InitMMapCachePool 用作需要通过文件持久化的时候，通过mmap方式进行更高效的文件i/o
// 通过内存mmap方式，把内存中的Chunk持久化到文件
// size: mmap方式一次分配的内存大小
// cnt:  每一次，预备多少块size大小的内存池
func InitMMapCachePool(dir string, size int, cnt int, errorfunc func(error), reloadfunc func([]*MMapCache)) error {
	os.MkdirAll(dir, os.ModePerm)
	DefPoolMMapCache = &PoolMMapCache{
		dir:        dir,
		template:   createMMapTemplate(size),
		recycleDur: time.Second,
		allocator:  make(chan *MMapCache),
		collector:  make(chan *MMapCache),
		errorfuc:   errorfunc,
		loadFlag:   false,
	}

	reload := DefPoolMMapCache.reloadCache()
	reloadfunc(reload)
	DefPoolMMapCache.mmapAllocLoop(cnt)
	for {
		if DefPoolMMapCache.loadFlag == false {
			<-time.After(time.Millisecond * 50)
		} else {
			break
		}
	}
	return nil
}

func createMMapTemplate(size int) []byte {
	template := make([]byte, size)
	for index := 0; index < size; index++ {
		template[index] = mmapInitByte
	}
	return template
}

func createMMapFile(file string, template []byte) error {
	return ioutil.WriteFile(file, template, 0666)
}

func (m *PoolMMapCache) reloadCache() []*MMapCache {
	for index := 0; ; index++ {
		filePath := path.Join(m.dir, fmt.Sprintf("%v.dat", index))
		_, err := os.Stat(filePath)
		if nil != err {
			break
		}
		os.Remove(filePath)
	}
	return nil
}

// Alloc 分配一个mmapcache
func (m *PoolMMapCache) Alloc() *MMapCache {
	return <-m.allocator
}

// Collect 回收一个mmapcache到缓存池
func (m *PoolMMapCache) Collect(mmcache *MMapCache) {
	m.collector <- mmcache
}

func (m *PoolMMapCache) preAllocMMapCache(fileid uint64) *MMapCache {
	filePath := path.Join(m.dir, fmt.Sprintf("%v.dat", fileid))

	createFlag := false
	fi, _ := os.Stat(filePath)
	if nil != fi {
		if int(fi.Size()) < len(m.template) {
			os.Remove(filePath)
			createFlag = true
		}
	} else {
		createFlag = true
	}

	if createFlag {
		err := createMMapFile(filePath, m.template)
		if nil != err {
			m.errorfuc(err)
			os.Exit(0)
		}
	}

	mmapCache, err := newMMapCache(filePath)
	if nil != err {
		m.errorfuc(err)
		os.Exit(0)
	}
	return mmapCache
}

func (m *PoolMMapCache) mmapAllocLoop(cnt int) {
	go func() {
		mmapIDSeq := uint64(0)
		q := list.New()
		for index := 0; index < cnt; index++ {
			q.PushBack(m.preAllocMMapCache(m.allocCounter))
			m.allocCounter++
		}
		m.loadFlag = true

		for {
			if q.Len() == 0 {
				q.PushBack(m.preAllocMMapCache(m.allocCounter))
				m.allocCounter++
			}
			e := q.Back()

			id := uint64(time.Now().Unix()) << 32
			id |= mmapIDSeq
			mmapIDSeq++
			e.Value.(*MMapCache).name(id)

			select {
			case b := <-m.collector:
				if q.Len() < cnt {
					b.recycle(m.template)
					q.PushBack(b)
					m.collectCounter++
				}
			case m.allocator <- e.Value.(*MMapCache):
				q.Remove(e)
			case <-time.After(m.recycleDur):
				if q.Len() > cnt {
					e := q.Back()
					e.Value.(*MMapCache).close(true)
					q.Remove(e)
					m.releaseCounter++
				} else {
					break
				}
			}
		}
	}()
}
