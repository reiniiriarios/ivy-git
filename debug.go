package main

import (
	"fmt"
	"runtime"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const STAT_INTERVAL = 10

var mem runtime.MemStats

func (a *App) statLoop() {
	for range time.Tick(time.Second * STAT_INTERVAL) {
		a.printStats()
	}
}

func (a *App) printStats() {
	runtime.ReadMemStats(&mem)
	wailsruntime.LogDebug(a.ctx, "mem.Alloc: "+byteCountIEC(mem.Alloc))
	wailsruntime.LogDebug(a.ctx, "mem.TotalAlloc: "+byteCountIEC(mem.TotalAlloc))
	wailsruntime.LogDebug(a.ctx, "mem.HeapAlloc: "+byteCountIEC(mem.HeapAlloc))
	wailsruntime.LogDebug(a.ctx, "mem.NumGC: "+fmt.Sprint(mem.NumGC))
}

func byteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
