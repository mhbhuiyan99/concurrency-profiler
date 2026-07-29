# concurrency-profiler

A Beego project that implements the same batch of API calls three different ways — sequential, `sync.WaitGroup`, and channels — then profiles each approach using Go's standard `runtime` and `runtime/pprof` tooling.

The goal is to compare execution time and allocation volume across the three concurrency strategies and report the results in a structured terminal output.

## Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Setup](#setup)
- [Architecture](#architecture)
- [External APIs](#external-apis)
- [Request Layer](#request-layer)
- [Concurrency Strategies](#concurrency-strategies)
- [Profiling](#profiling)
- [Timing Explanation](#timing-explanation)
- [Performance Comparison Logic](#performance-comparison-logic)
- [Why Results Differ Between Runs](#why-results-differ-between-runs)
- [Sample Output](#sample-output)
- [Notes](#notes)
- [Expected Outcome](#expected-outcome)

## Overview

`GET /test-concurrency` runs the same batch of 12 API calls three times, once per strategy, and returns:

- Execution time for each strategy
- Allocation volume (bytes allocated) for each strategy
- Controller-level memory statistics (`Alloc`, `TotalAlloc`, `Sys`, `NumGC`) before and after the full run
- A CPU profile written to `profiles/cpu.prof`, captured across the entire controller execution
- A comparison summary: fastest method, highest/lowest allocation, and a heuristic "most efficient" method

Goroutines are used by the `WaitGroup` and `Channel` implementations to run requests concurrently, but the project does **not** report a phase-level goroutine count. The assignment lists goroutine count as "optional but preferred" — it was intentionally left out because `runtime.NumGoroutine()` reflects the whole Go process rather than the phase currently running, so a per-phase reading would be noisy and not a reliable comparison metric. Profiling is scoped to execution time and allocation volume instead.

## Project Structure

```
concurrency-profiler/
├── controllers/
│   └── concurrency_controller.go
├── concurrency/
│   ├── channel.go
│   ├── fetch.go
│   ├── sequential.go
│   ├── types.go
│   ├── urls.go
│   └── waitgroup.go
├── profiling/
│   ├── comparison.go
│   ├── cpu.go
│   ├── memory.go
│   ├── phase.go
│   └── types.go
├── requests/
│   ├── client.go
│   └── request.go
├── routers/
│   └── router.go
├── utils/
│   ├── report.go
│   ├── timer.go
│   └── visual.go
├── conf/
│   ├── app.conf            # local configuration, not committed
│   └── app.conf.example    # configuration template
├── profiles/
│   └── cpu.prof             # generated on each run, not committed
├── main.go
├── go.mod
└── go.sum
```

### Layer responsibilities

```
Controller        → HTTP request handling, orchestrates the benchmark run
Concurrency layer → the three execution strategies (Sequential / WaitGroup / Channel)
Request layer     → all external HTTP communication (auth, headers, timeouts)
Profiling layer   → CPU profiling + memory snapshots
Utils layer       → timing helper + terminal report / visualization
```

Request flow:

```
Controller
   ↓
Concurrency layer
   ↓
Request layer
   ↓
External APIs
```

Profiling and reporting sit alongside this flow rather than inside it — `profiling.ProfilePhase` wraps each concurrency call, and `utils` formats what comes out:

```
Profiling layer → measures execution time + allocation volume per phase
Utils layer     → prints the terminal report from those measurements
```

Controllers never call external APIs directly; all HTTP calls go through `requests/`. This separation is what lets the same `FetchAPI` function be reused unchanged by all three concurrency strategies.

## Setup

Requires Go 1.25+.

```bash
git clone https://github.com/mhbhuiyan99/concurrency-profiler.git
cd concurrency-profiler
cp conf/app.conf.example conf/app.conf
```

Fill in `conf/app.conf` with real values:

```ini
appname = concurrency-profiler
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true

base_url = <your_base_url>
username = <your_username>
password = <your_password>
api_key  = <your_api_key>
```

`conf/app.conf` holds live credentials and is intentionally excluded from version control (see `.gitignore`); only `conf/app.conf.example` is committed.

Run the server:

```bash
go run main.go
```

Trigger the benchmark:

```bash
curl http://localhost:8080/test-concurrency
```

## Architecture

```
TestConcurrency()
│
├── start controller timer
├── Start CPU profiling
├── memory before
│
├── ProfilePhase("Sequential")
├── ProfilePhase("WaitGroup")
├── ProfilePhase("Channel")
│
├── memory after
├── ComparePhases()
├── Stop CPU profiling
├── stop controller timer
│
├── terminal report
└── JSON response
```

Each phase follows the same request path down to the HTTP layer:

```
ProfilePhase()
   │
   ├── MeasureExecution()
   │      │
   │      └── Run<Strategy>()
   │             │
   │             └── FetchAPI()
   │                    │
   │                    ├── NewGETRequest()
   │                    ├── setDefaultHeaders()
   │                    └── DoRequest()
```

## External APIs

The benchmark calls 12 Presto category-detail endpoints, grouped roughly around three regions:

- `usa:hawaii`
- `north-america`
- `usa:texas`

The list is defined once in `concurrency.APIURLs` (`concurrency/urls.go`) and reused by all three strategies, so there is a single source of truth for the workload.

**"Total APIs Called: 12" refers to the workload size per execution method, not the total number of network requests made by the benchmark.** Each phase — Sequential, WaitGroup, and Channel — calls all 12 endpoints independently:

```
Sequential → 12 API requests
WaitGroup  → 12 API requests
Channel    → 12 API requests
-----------------------------
Total requests across one full benchmark run = 36
```

## Request Layer

`requests/client.go` defines a single shared `http.Client` with a 10-second timeout, reused across every call to avoid re-creating a client per request.

`requests/request.go` builds each request:

- `NewGETRequest` creates the request and applies default headers.
- `setDefaultHeaders` reads `username`, `password`, and `api_key` from `conf/app.conf`, applies HTTP Basic Auth, and sets the remaining headers (`Accept`, `Content-Type`, `X-Requested-With`, `User-Agent`, `Origin`, `x-api-key`).
- `DoRequest` executes the request, validates the status code, and reads the response body.

`concurrency.FetchAPI` is the single entry point used by all three strategies — it calls `NewGETRequest` and `DoRequest` and wraps the outcome (or error) into an `APIResult`.

## Concurrency Strategies

### 1. Sequential

`concurrency/sequential.go` loops over the URL list and calls `FetchAPI` one at a time, appending each result before moving to the next.

### 2. WaitGroup

`concurrency/waitgroup.go` launches one goroutine per URL, writes each result into a pre-sized slice by index (avoiding a shared-append race), and blocks on `sync.WaitGroup.Wait()` until all goroutines finish.

### 3. Channel

`concurrency/channel.go` launches one goroutine per URL, each sending its result into a channel buffered to the known URL count. A separate goroutine closes the channel once `wg.Wait()` returns, and the main goroutine drains results with a `range` loop.

## Profiling

### Controller-level

Wraps the entire `/test-concurrency` execution:

```
Controller: /test-concurrency
│
├── Start CPU profiling
├── Memory stats before (Alloc, TotalAlloc, Sys, NumGC)
├── Sequential phase
├── WaitGroup phase
├── Channel phase
├── Memory stats after (Alloc, TotalAlloc, Sys, NumGC)
└── Stop CPU profiling
```

### Phase-level

Each phase records memory statistics immediately before and after it runs:

```
ProfilePhase(name, fn)
│
├── memory before  (runtime.ReadMemStats)
├── MeasureExecution(fn)   — timed call to the strategy
└── memory after   (runtime.ReadMemStats)
```

`profiling.CalculateMemoryUsed` takes the difference between `TotalAlloc` after and before. This is **allocation volume** — total bytes allocated during the phase — not live/retained memory. It will not decrease even if the garbage collector frees memory during the phase, since `TotalAlloc` only ever increases.

### CPU profiling

`profiling/cpu.go` wraps `pprof.StartCPUProfile()` / `pprof.StopCPUProfile()` around the **full controller execution** — all three phases plus comparison and reporting — and writes the result to `profiles/cpu.prof`.

There is no separate CPU profile per phase. The profile file reflects the whole `/test-concurrency` request, not Sequential, WaitGroup, or Channel individually. Inspect it with `go tool pprof profiles/cpu.prof` if you want a breakdown, but that breakdown will span all three strategies together, not isolate one.

### Memory / allocation profiling

Two separate memory measurements exist and should not be conflated:

- **Controller-level**: one `Alloc`/`TotalAlloc`/`Sys`/`NumGC` snapshot before all three phases run, and one after — a single before/after pair for the whole request.
- **Phase-level**: a before/after `TotalAlloc` pair captured around each individual phase, reduced to a single "bytes allocated" number per phase.

## Timing Explanation

- **Phase-level time**: `utils.MeasureExecution` starts a timer immediately before calling the strategy function (`RunSequential`, `RunWaitGroup`, or `RunChannel`) and stops it right after it returns.
- **Controller-level time**: `TestConcurrency` times the whole request handler, from before CPU profiling starts to after all three phases and the comparison step complete.

Controller total time will always be somewhat larger than the sum of the three phase times, since it also includes profiling overhead, comparison logic, and report generation.

## Performance Comparison Logic

`profiling.ComparePhases` produces:

- **Fastest method** — the strategy with the lowest execution time.
- **Highest / lowest allocation** — the strategy with the highest / lowest `TotalAlloc` delta.
- **Performance gain** — percentage difference in execution time relative to a baseline (`WaitGroup vs Sequential`, `Channel vs Sequential`, `Channel vs WaitGroup`).
- **Most efficient method** — a project-specific heuristic, not a universal definition of efficiency. A method is "most efficient" if it has *both* the lowest execution time and lowest allocation volume among the three. If no single method dominates both metrics, the fastest method is used as the fallback.

## Why Results Differ Between Runs

This benchmark calls real external APIs, so results vary from run to run and WaitGroup and Channel do not have a fixed winner. Contributing factors include:

- Network latency and remote server response time on each request
- Whether the underlying HTTP connection is reused or freshly established (TLS handshake cost)
- Go's goroutine scheduling for a given run
- Channel synchronization overhead vs. WaitGroup synchronization overhead
- External API or upstream caching behavior

WaitGroup is not inherently always faster than Channel, and Channel is not inherently always faster than WaitGroup — each run measures what happened during that particular set of 36 requests, not a fixed property of the language feature. Run the benchmark multiple times if you want a representative picture rather than trusting a single run.

## Sample Output

Values below are illustrative, not real captured output — your numbers will vary between runs for the reasons described above.

```
================ API PERFORMANCE TEST ================

Total APIs Called: 12

[ Sequential Execution ]
-------------------------------------
Time Taken : 3182.44 ms

[ WaitGroup Execution ]
-------------------------------------
Time Taken : 941.07 ms

[ Channel Execution ]
-------------------------------------
Time Taken : 863.21 ms

================= COMPARISON =================

Execution Time
-------------------------------------
Sequential | ██████████████████████████████ 3.18 s
WaitGroup  | █████████                       0.94 s
Channel    | ████████                        0.86 s

Performance Gain
-------------------------------------
WaitGroup vs Sequential : 70.43% faster
Channel vs Sequential   : 72.87% faster
Channel vs WaitGroup    : 8.27% faster

Allocation Volume
-------------------------------------
Highest Allocation : Sequential
Lowest Allocation  : Channel

================ PROFILING REPORT ================

Controller : /test-concurrency

Total Execution Time : 5001.32 ms

--------------- MEMORY STATS ----------------

Before Execution
Alloc      : 2.14 MB
TotalAlloc : 2.14 MB
Sys        : 6.53 MB
NumGC      : 2

After Execution
Alloc      : 3.79 MB
TotalAlloc : 5.02 MB
Sys        : 8.20 MB
NumGC      : 4

--------------- PHASE PROFILING ----------------

Sequential
Time Taken      : 3182.44 ms
Bytes Allocated : +1.48 MB

WaitGroup
Time Taken      : 941.07 ms
Bytes Allocated : +1.05 MB

Channel
Time Taken      : 863.21 ms
Bytes Allocated : +0.97 MB

--------------- CPU PROFILE ----------------

CPU Profiling Started...
CPU Profiling Stopped.

--------------- SUMMARY ----------------

Fastest Method           : Channel
Highest Allocation       : Sequential
Most Efficient Method    : Channel

====================================================
```

## Notes

- Goroutines are used by the WaitGroup and Channel implementations, but phase-level goroutine counts are intentionally not reported (see [Overview](#overview)).
- CPU profiling covers the whole controller execution, not each phase individually (see [Profiling](#profiling)).
- Phase-level memory numbers represent allocation volume (`TotalAlloc` delta), not live memory usage.
- "Total APIs Called: 12" is the per-phase workload; the full benchmark makes 36 external requests.

## Expected Outcome

This project demonstrates:

- Core Go concurrency patterns: goroutines, channels (buffered), and `sync.WaitGroup`
- Sequential vs. concurrent execution of the same workload
- Execution-time measurement via `time.Since`
- Runtime memory statistics via `runtime.ReadMemStats`
- Allocation-volume comparison across execution strategies
- CPU profiling via `runtime/pprof`
- Terminal-based visualization and structured performance reporting
- Interpreting profiling output to identify the fastest and most efficient execution strategy for a given workload, while understanding why results vary between runs against a live external API