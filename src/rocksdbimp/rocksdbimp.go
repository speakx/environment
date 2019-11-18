package rocksdbimp

import (
	"github.com/tecbot/gorocksdb"
)

// RocksdbImp 包装类
type RocksdbImp struct {
	db       *gorocksdb.DB
	writeOpt *gorocksdb.WriteOptions
	readOpt  *gorocksdb.ReadOptions
}

// NewRocksdbImp new RocksdbImp
func NewRocksdbImp() *RocksdbImp {
	r := &RocksdbImp{
		writeOpt: getDefWriteOptions(),
		readOpt:  getDefReadOptions(),
	}
	return r
}

// OpenDB open rocksdb
func (r *RocksdbImp) OpenDB(dbPath string) error {
	db, err := gorocksdb.OpenDb(getDefOptions(), dbPath)
	if err != nil {
		db.Close() // 官方example中对err的case也会调用Close方法
		return err
	}

	r.db = db
	return nil
}

// GetDB return rocksdb
func (r *RocksdbImp) GetDB() *gorocksdb.DB {
	return r.db
}

// Close close
func (r *RocksdbImp) Close() {
	r.db.Close()
}

// Put rocksdb.Put
func (r *RocksdbImp) Put(key, value []byte) error {
	return r.db.Put(r.writeOpt, key, value)
}

// Get rocksdb.Get
func (r *RocksdbImp) Get(key []byte) (*gorocksdb.Slice, error) {
	return r.db.Get(r.readOpt, key)
}

// GetWriteOptions WriteOptions
func (r *RocksdbImp) GetWriteOptions() *gorocksdb.WriteOptions {
	return r.writeOpt
}

// GetReadOptions ReadOptions
func (r *RocksdbImp) GetReadOptions() *gorocksdb.ReadOptions {
	return r.readOpt
}

func getDefOptions() *gorocksdb.Options {
	opt := gorocksdb.NewDefaultOptions()
	opt.SetCreateIfMissing(true)
	opt.SetCompression(gorocksdb.SnappyCompression)
	opt.SetMaxBackgroundFlushes(10)
	opt.SetCreateIfMissingColumnFamilies(true)
	opt.SetWriteBufferSize(1 << 30)
	opt.SetMaxWriteBufferNumber(1 << 32)
	opt.SetArenaBlockSize(1 << 28)
	opt.SetKeepLogFileNum(3)
	opt.SetInfoLogLevel(2)
	opt.SetMaxOpenFiles(32768)
	return opt
}

func getDefWriteOptions() *gorocksdb.WriteOptions {
	opt := gorocksdb.NewDefaultWriteOptions()
	return opt
}

func getDefReadOptions() *gorocksdb.ReadOptions {
	opt := gorocksdb.NewDefaultReadOptions()
	// opt.SetPrefixExtractor(gorocksdb.NewFixedPrefixTransform(3))
	return opt
}
