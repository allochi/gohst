package main

import (
	// "fmt"
	"runtime"
	// "time"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	spew.Dump(m)

}
