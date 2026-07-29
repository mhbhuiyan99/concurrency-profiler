package profiling

import "runtime"

// GetMemoryStats collects the current memory statistics.
//
// Responsibilities:
//   - Read the current runtime memory statistics.
//   - Extract the required memory metrics.
//   - Return the collected memory information.
func GetMemoryStats() MemoryStats {

	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	return MemoryStats{
		Alloc:      mem.Alloc,
		TotalAlloc: mem.TotalAlloc,
		Sys:        mem.Sys,
		NumGC:      mem.NumGC,
	}
}