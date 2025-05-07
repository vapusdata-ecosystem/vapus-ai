package dmutils

import (
	"fmt"
	"runtime"
)

func LogMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc <==> %.2f MiB", float64(m.Alloc)/1024/1024)
}

func CleanPointers(agent any) {
	LogMemoryUsage()
	agent = nil
	runtime.GC()
	LogMemoryUsage()
}
