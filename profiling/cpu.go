package profiling

import (
	"fmt"
	"os"
	"runtime/pprof"
)

// StartCPUProfile starts CPU profiling and writes the profile to a file.
//
// Responsibilities:
//   - Create the CPU profile file.
//   - Start Go CPU profiling.
//   - Return the profile file for later cleanup.
func StartCPUProfile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("create CPU profile file: %w", err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		file.Close()
		return nil, fmt.Errorf("start CPU profiling: %w", err)
	}

	return file, nil
}

// StopCPUProfile stops the active CPU profiler.
//
// Responsibilities:
//   - Stop Go CPU profiling.
//   - Close the CPU profile file.
func StopCPUProfile(file *os.File) error {
	pprof.StopCPUProfile()
	
	if err := file.Close(); err != nil {
		return fmt.Errorf("close CPU profile file: %w", err)
	}

	return nil
}