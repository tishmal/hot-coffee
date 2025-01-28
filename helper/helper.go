package helper

import (
	"fmt"
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

func PrintUsage() {
	fmt.Println(`$ ./hot-coffee --help
Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`)
}
