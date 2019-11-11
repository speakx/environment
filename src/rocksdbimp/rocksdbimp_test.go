package rocksdbimp

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/tecbot/gorocksdb"
)

var db *RocksdbImp
var writeOpt *gorocksdb.WriteOptions
var key = []byte("Hello")
var val = []byte("World")

func TestMmapRocksdb(t *testing.T) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	db = NewRocksdbImp()
	err := db.OpenDB(dir)
	if nil != err {
		t.Errorf("failed to open rocksdb, err:%v", err)
	}
	writeOpt = gorocksdb.NewDefaultWriteOptions()
}
func Benchmark_rocksdb_kv_put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := db.GetDB().Put(writeOpt, key, val)
		if nil != err {
			b.Errorf("failed to rocksdb.put err:%v", err)
		}
	}
}

func Benchmark_rocksdb_kv_put_goroutine(b *testing.B) {
	routineCount := 100
	goroutine := b.N / routineCount
	mCount := 0
	if b.N < routineCount {
		goroutine = 1
		mCount = b.N
	} else {
		if goroutine > 200 {
			goroutine = 200
			routineCount = b.N / goroutine
		}
		if goroutine*routineCount < b.N {
			goroutine++
			mCount = goroutine*routineCount - b.N
		}
	}

	b.Logf("n:%v goroutine:%v c:%v m:%v", b.N, goroutine, routineCount, mCount)
	wait := &sync.WaitGroup{}
	wait.Add(goroutine)
	for index := 0; index < goroutine; index++ {
		go func(idx int) {
			writeCnt := routineCount
			if idx == goroutine-1 {
				writeCnt = mCount
			}
			for i := 0; i < writeCnt; i++ {
				err := db.GetDB().Put(writeOpt, key, val)
				if nil != err {
					b.Errorf("failed to rocksdb.put err:%v", err)
				}
			}
			wait.Done()
		}(index)
	}
	wait.Wait()
}
