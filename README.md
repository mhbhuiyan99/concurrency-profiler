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