# golang-calculator
Simple GO web-app with calculator
---
## 🛫 Available Commands
Use the Makefile to run common commands:

```bash
make help           # Show all available commands
make run            # Start the application
make test           # Run unit tests
make coverage       # Run tests with coverage report
make coverage-html  # Run tests and open coverage in HTML browser
make build          # Build the application binary
```

---
## 🚀 Quick Start
Start the app:
```bash
make run
```

---
## 🧪 Running Unit Tests
Run all tests:
```bash
make test
```

Run tests with coverage report:
```bash
make coverage
```

View coverage in HTML format:
```bash
make coverage-html
```

Run tests for a specific directory:
```bash
go test -v ./internal/domain
```
