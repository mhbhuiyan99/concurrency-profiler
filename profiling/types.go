package profiling

// MemoryStats stores memory usage statistics.
//
// Responsibilities:
//   - Store allocated memory.
//   - Store total allocated memory.
//   - Store system memory.
//   - Store garbage collection count.
type MemoryStats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}