```
Controller

│

▼

MeasureExecution()

│

├── start timer

│

├── fn()

│      │
│      ▼
│   RunSequential()
│      │
│      ├── URL 1
│      │     │
│      │     ▼
│      │ NewGETRequest()
│      │
│      │ DoRequest()
│      │
│      ├── URL 2
│      │
│      ├── URL 3
│      │
│      └── ...
│
│
├── stop timer
│
▼

PhaseResult

│

▼

Controller
```