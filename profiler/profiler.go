package profiler

import (
	"runtime"
	"sync"
	"time"
)

func New() *Profiler {
	profiler := &Profiler{
		start: time.Now(),
	}

	profiler.loadCurrentMemory()

	go profiler.poller()

	return profiler
}

type Profiler struct {
	peakMemory    uint64
	currentMemory uint64
	start         time.Time
	mu            sync.Mutex
}

func (profiler *Profiler) GetCurrentMemory() uint64 {
	profiler.mu.Lock()
	defer profiler.mu.Unlock()

	return profiler.currentMemory
}

func (profiler *Profiler) GetPeakMemory() uint64 {
	profiler.mu.Lock()
	defer profiler.mu.Unlock()

	return profiler.peakMemory
}

func (profiler *Profiler) GetDuration() time.Duration {
	profiler.mu.Lock()
	defer profiler.mu.Unlock()

	return time.Now().Sub(profiler.start)
}

func (profiler *Profiler) loadCurrentMemory() {
	profiler.mu.Lock()
	defer profiler.mu.Unlock()

	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)

	profiler.currentMemory = memory.Sys
	if profiler.currentMemory > profiler.peakMemory {
		profiler.peakMemory = profiler.currentMemory
	}
}

func (profiler *Profiler) poller() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		profiler.loadCurrentMemory()
	}
}
