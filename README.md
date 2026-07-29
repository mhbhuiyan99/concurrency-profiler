## Architecture
```
Controller
    │
    ├── StartCPUProfile()
    │
    ├── GetMemoryStats()
    │
    ├── ProfilePhase()
    │      │
    │      └── MeasureExecution()
    │             │
    │             └── RunSequential()
    │
    ├── ProfilePhase()
    │      │
    │      └── MeasureExecution()
    │             │
    │             └── RunWaitGroup()
    │
    ├── ProfilePhase()
    │      │
    │      └── MeasureExecution()
    │             │
    │             └── RunChannel()
    │
    ├── GetMemoryStats()
    │
    └── StopCPUProfile()
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