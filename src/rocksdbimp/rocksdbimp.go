package rocksdbimp

import (
	"github.com/tecbot/gorocksdb"
)

// RocksdbImp 包装类
type RocksdbImp struct {
	db *gorocksdb.DB
}

// NewRocksdbImp new RocksdbImp
func NewRocksdbImp() *RocksdbImp {
	r := &RocksdbImp{}
	return r
}

// OpenDB open rocksdb
func (r *RocksdbImp) OpenDB(dbPath string) error {
	options := r.getOptions()
	options.SetCreateIfMissing(true)

	db, err := gorocksdb.OpenDb(options, dbPath)
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

func (r *RocksdbImp) getOptions() *gorocksdb.Options {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetCompression(gorocksdb.SnappyCompression)
	opts.SetMaxBackgroundFlushes(10)
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetWriteBufferSize(1 << 30)
	opts.SetMaxWriteBufferNumber(1 << 32)
	opts.SetArenaBlockSize(1 << 28)
	opts.SetKeepLogFileNum(3)
	opts.SetInfoLogLevel(2)
	opts.SetMaxOpenFiles(32768)
	return opts
}
