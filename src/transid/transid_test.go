package transid

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func Benchmark_uuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.NewV1().String()
	}
}
