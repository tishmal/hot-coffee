package helper

import (
	"sync/atomic"
)

var counter uint64

// инкрементит безопасно id
func GenerateID() uint64 {
	return atomic.AddUint64(&counter, 1)
}
