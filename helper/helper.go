package helper

import (
	"sync/atomic"
)

var counter uint64

// инкрементит безопасно id
func GenerateID() uint64 {
	return atomic.AddUint64(&counter, 1)
}

func ErrorResponseInJSON() {
	// реализовать для повсеместного применения в коде. Возврат ошибок в формате JSON
}
