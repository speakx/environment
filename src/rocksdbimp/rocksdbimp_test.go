package rocksdbimp

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/tecbot/gorocksdb"
)

var db *RocksdbImp
var writeOpt *gorocksdb.WriteOptions
var readOpt *gorocksdb.ReadOptions
var itemCounter = 0
var keyFMT = "Hello-%08X"
var valFMT = "World-%v"
var r *rand.Rand
var mapItemCounter = 0
var mapKV map[string]string

func TestRocksdb(t *testing.T) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	db = NewRocksdbImp()
	err := db.OpenDB(dir)
	if nil != err {
		t.Errorf("failed to open rocksdb, err:%v", err)
	}
	r = rand.New(rand.NewSource(time.Now().Unix()))
	mapKV = make(map[string]string)

	db.Put([]byte("Hello"), []byte("World"))
	val, err := db.Get([]byte("HelloX"))
	t.Logf("val:%v exist:%v err:%v", val, val.Exists(), err)
	i, err := strconv.ParseUint(string(val.Data()), 10, 64)
	t.Logf("empty ParseUint => %v err:%v", i, err)
}

func TestRocksdb_range(t *testing.T) {
	for index := 0; index < 50; index++ {
		err := db.GetDB().Put(db.GetWriteOptions(),
			[]byte(fmt.Sprintf(keyFMT, itemCounter)),
			[]byte(fmt.Sprintf(valFMT, itemCounter)))
		itemCounter++
		if nil != err {
			t.Errorf("failed to rocksdb.put err:%v", err)
			return
		}
	}

	it := db.GetDB().NewIterator(db.GetReadOptions())
	defer it.Close()

	it.Seek([]byte(fmt.Sprintf(keyFMT, 10)))
	for ; it.Valid(); it.Next() {
		t.Logf("Key: %v Value: %v", string(it.Key().Data()), string(it.Value().Data()))
	}
}

func Benchmark_put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := db.Put(
			[]byte(fmt.Sprintf(keyFMT, itemCounter)),
			[]byte(fmt.Sprintf(valFMT, itemCounter)))
		itemCounter++
		if nil != err {
			b.Errorf("failed to rocksdb.put err:%v", err)
		}
	}
}

func Benchmark_get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := db.Get([]byte(fmt.Sprintf(keyFMT, r.Int()%itemCounter)))
		if nil != err {
			b.Errorf("failed to rocksdb.get err:%v", err)
		}
	}
}

func Benchmark_range(b *testing.B) {
	b.Logf("KeyRange [%v, %v)", fmt.Sprintf(keyFMT, 0), fmt.Sprintf(keyFMT, itemCounter))

	for i := 0; i < b.N; i++ {
		it := db.GetDB().NewIterator(db.GetReadOptions())
		defer it.Close()

		it.Seek([]byte(fmt.Sprintf(keyFMT, r.Int()%itemCounter)))
		for n := 0; it.Valid() && n < 50; it.Next() {
			n++
		}
	}
}

func Benchmark_map_put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mapKV[fmt.Sprintf(keyFMT, mapItemCounter)] = fmt.Sprintf(valFMT, mapItemCounter)
		mapItemCounter++
	}
}

func Benchmark_map_get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf(keyFMT, r.Int()%itemCounter)
		_, ok := mapKV[key]
		if false == ok {
			b.Errorf("%v not found", key)
		}
	}
}

func Benchmark_ParseUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseUint("9223372036854775807", 10, 64)
	}
}

func Benchmark_Atoi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Atoi("9223372036854775807")
	}
}

func Benchmark_FmtInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", 9223372036854775807)
	}
}

func Benchmark_Itoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(9223372036854775807)
	}
}

func Benchmark_FormatUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.FormatUint(9223372036854775807, 10)
	}
}

// func Benchmark_rocksdb_kv_put_goroutine(b *testing.B) {
// 	routineCount := 100
// 	goroutine := b.N / routineCount
// 	mCount := 0
// 	if b.N < routineCount {
// 		goroutine = 1
// 		mCount = b.N
// 	} else {
// 		if goroutine > 200 {
// 			goroutine = 200
// 			routineCount = b.N / goroutine
// 		}
// 		if goroutine*routineCount < b.N {
// 			goroutine++
// 			mCount = goroutine*routineCount - b.N
// 		}
// 	}

// 	b.Logf("n:%v goroutine:%v c:%v m:%v", b.N, goroutine, routineCount, mCount)
// 	wait := &sync.WaitGroup{}
// 	wait.Add(goroutine)
// 	for index := 0; index < goroutine; index++ {
// 		go func(idx int) {
// 			writeCnt := routineCount
// 			if idx == goroutine-1 {
// 				writeCnt = mCount
// 			}
// 			for i := 0; i < writeCnt; i++ {
// 				err := db.GetDB().Put(writeOpt, key, val)
// 				if nil != err {
// 					b.Errorf("failed to rocksdb.put err:%v", err)
// 				}
// 			}
// 			wait.Done()
// 		}(index)
// 	}
// 	wait.Wait()
// }
