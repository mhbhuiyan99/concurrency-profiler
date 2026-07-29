## Architecture
```
Controller
    │
    ├── Start CPU profiling
    │
    ├── Memory Before
    │
    ├── ProfilePhase("Sequential")
    │       └── MeasureExecution()
    │
    ├── ProfilePhase("WaitGroup")
    │       └── MeasureExecution()
    │
    ├── ProfilePhase("Channel")
    │       └── MeasureExecution()
    │
    ├── Memory After
    │
    └── ComparePhases()
            │
            ├── Fastest
            ├── Performance gain
            ├── Memory comparison
            ├── Goroutine comparison
            └── Most efficient
```

### Sequential
```
Client
   │
   ▼
Router
   │
   ▼
ConcurrencyController.TestConcurrency()
   │
   ▼
MeasureExecution()
   │
   ├── Start Timer
   │
   ├── RunSequential()
   │       │
   │       ├── For each URL
   │       │       │
   │       │       ▼
   │       │   FetchAPI()
   │       │       │
   │       │       ▼
   │       │   NewGETRequest()
   │       │       │
   │       │       ▼
   │       │   setDefaultHeaders()
   │       │       │
   │       │       ▼
   │       │   DoRequest()
   │       │       │
   │       │       ▼
   │       │   httpClient.Do()
   │       │       │
   │       │       ▼
   │       │   APIResult
   │       │
   │       └── Return []APIResult
   │
   ├── Stop Timer
   └── Build PhaseResult
   │
   ▼
ServeJSON()
```

### WaitGroup

```
RunWaitGroup()

        │
        ▼

Create WaitGroup

        │
        ▼

Loop URLs

        │
        ├── Add(1)
        │
        ├── Goroutine 1
        │        │
        │        ▼
        │    FetchAPI()
        │
        ├── Goroutine 2
        │
        ├── Goroutine 3
        │
        └── ...

        │
        ▼

Wait()

        │
        ▼

Return []APIResult
```

### Channel

```
RunChannel()

        │
        ▼

Create Channel

        │
        ▼

Loop URLs

        │
        ├── Goroutine 1
        │        │
        │        ▼
        │   FetchAPI()
        │        │
        │        ▼
        │   channel <- APIResult
        │
        ├── Goroutine 2
        │
        ├── Goroutine 3
        │
        └── ...

        │
        ▼

Main Goroutine

        │
        ▼

Receive 12 Results

        │
        ▼

Return []APIResult
```

## Profiling

```
                    Controller
                        │
                        ▼
              ProfilePhase(...)
                        │
            ┌───────────┴───────────┐
            │                       │
            ▼                       ▼
    Memory/Goroutine          MeasureExecution()
       profiling                   │
            │                      │
            │                      ▼
            │                 RunSequential()
            │
            └───────────┬───────────┘
                        │
                        ▼
                 ProfileResult
```

```
ProfilePhase()
    │
    ├── GetMemoryStats()
    │
    ├── NumGoroutine()
    │
    ├── MeasureExecution()
    │       │
    │       ├── start timer
    │       │
    │       ├── fn()
    │       │      │
    │       │      └── RunSequential()
    │       │
    │       └── stop timer
    │
    ├── GetMemoryStats()
    │
    └── NumGoroutine()
```

### Controller-level
Controller-level profiling measures the overall `/test-concurrency` execution. It captures memory statistics before and after all three concurrency phases and runs CPU profiling across the complete controller execution.
```
Controller: /test-concurrency
│
├── Start CPU profiling
│
├── Memory statistics before execution
│      ├── Alloc
│      ├── TotalAlloc
│      ├── Sys
│      └── NumGC
│
├── Sequential Phase
│
├── WaitGroup Phase
│
├── Channel Phase
│
├── Memory statistics after execution
│      ├── Alloc
│      ├── TotalAlloc
│      ├── Sys
│      └── NumGC
│
└── Stop CPU profiling
```

### Phase-level

Phase-level goes inside that controller and measures each execution strategy separately.
```
Controller
│
├── Sequential
│      ├── memory before
│      ├── goroutines before
│      ├── execution time
│      ├── memory after
│      └── goroutines after
│
├── WaitGroup
│      ├── memory before
│      ├── goroutines before
│      ├── execution time
│      ├── memory after
│      └── goroutines after
│
└── Channel
       ├── memory before
       ├── goroutines before
       ├── execution time
       ├── memory after
       └── goroutines after
```
