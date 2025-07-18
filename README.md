# SEQUEX

**SEQUEX** is a centralized, event-driven trading system designed for high-performance and modularity. It provides a clean architecture for handling trading logic, messaging, infrastructure, and configuration in a scalable and maintainable way.

---

## 📦 Project Structure
```
.
├── cmd             # Application entry point
├── config          # Configuration files
├── env             # Environment definitions
├── go.mod
├── go.sum
├── internal        # Core logic
├── Makefile        # Automation command
├── pkg             # Shared utilities
└── README.md
```

## ⚙️ Make command


Build binaries to bin/
```
make build
```

Run
```
./bin/sequex-linux-x86 -c ./config/config_sequex.json
```
or
```
go run cmd/sequex/server.go -c ./config/config_sequex.json
```

Test
```
make test
```