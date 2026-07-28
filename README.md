```
Client (Browser / Postman)

        │
        ▼

Router

        │
        ▼

ConcurrencyController.TestConcurrency()

        │
        ▼

MeasureExecution("Sequential", RunSequential)

        │
        ├── Start Timer
        │
        ├── Execute fn()
        │
        ▼

RunSequential()

        │
        ├── Loop through all API URLs
        │
        ├──────────────┐
        │              │
        ▼              ▼

    FetchAPI(URL1)  FetchAPI(URL2) ... FetchAPI(URL12)
        │
        ▼

NewGETRequest()

        │
        ▼

setDefaultHeaders()

        │
        ▼

DoRequest()

        │
        ▼

httpClient.Do()

        │
        ▼

API Response

        │
        ▼

APIResult

        │
        ▼

Return []APIResult

        │
        ▼

MeasureExecution()

        │
        ├── Stop Timer
        ├── Calculate Execution Time
        └── Build PhaseResult

        │
        ▼

Return PhaseResult

        │
        ▼

ConcurrencyController

        │
        ▼

ServeJSON()
```