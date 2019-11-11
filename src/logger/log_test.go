package logger

import (
	"testing"
)

func Benchmark_rocksdb_kv_put(b *testing.B) {
	InitLogger("./test.log", false, "info")

	for i := 0; i < b.N; i++ {
		Info("Hello ", "world")
	}
}
