# Exchange Connectivity

## Purpose

The `exchange` module is a **collection of crypto-exchange connectivity implementations** for the `pkg`, designed to provide **reliable and scalable exchange access** for trading systems.

- The goal is **not** to strictly centralize all exchange APIs under a single unified abstraction, but to **customize implementations to each exchange's specific behavior while centralizing only minimal APIs and reusable patterns**.
- For **future development**, the `exchange` module:
  - **Should not have any exchange-specific module dependencies across exchanges**.
  - Should only rely on **low-level modules** such as `http`, `websocket`, and (potentially) `FIX` for communication.
  - Should remain composable, testable, and cleanly separated per exchange.

---

## Design

Each **exchange** is implemented as its **own submodule under the `exchange/` folder**:

- **API differences define submodules**:
  - For example, `binance` (spot) and `binancefutures` are **separate modules** due to different API behaviors.
  - Conversely, `binance` and its **testnet** are treated as the same module, with the environment determined by configuration.
- Some exchanges, such as Bybit (which uses unified accounts for PERP and SPOT), will implement **both PERP and SPOT under a single folder**, while retaining clear separation within the implementation.
- **Users are expected to understand exchange-specific nuances** before referencing the codebase for integration or extension.

---

## Implementation Structure

Each exchange implementation should cover, **at minimum**, the following structure (from low-level to high-level):

- `config.go` — Configures API parameters (e.g., `API_KEY`, `API_SECRET`, `BASE_URL`).
- `utils.go` — Utilities such as signing HTTP headers or timestamp generation.
- `const.go` — Constants such as base URLs, endpoint paths, and header keys, enums.
- `error.go` — Exchange-specific API error code parsing and mapping.
- `request.go` — Low-level HTTP requests, session management, and signed/unsigned requests.
- `models.go` — Models for HTTP requests/responses, enabling `json.Marshal` and `json.Unmarshal` to Go types.
- `client.go` — High-level wrapper over `request.go` with API-specific logic.
- `client_test.go` — Unit tests for `client.go`.
- `ws_models.go` — Models for WebSocket requests/responses for serialization/deserialization.
- `ws.go` — Low-level WebSocket connection handling:
  - Establish connection
  - Subscribe/unsubscribe logic
  - Retry, reconnection, and error handling
- `ws_client.go` — High-level WebSocket API:
  - Subscription interfaces
  - Event dispatching
  - Topic-specific data handlers
- `utils.go` — Miscellaneous helper function can be re-use under the corresponding exchanges

---

## Testing Guidelines

- After implementing each **high-level API** (`client.go`, `ws_client.go`), add **unit tests** to cover the **happy path and critical edge cases**.
- **Do not use mocks for unit tests**:
  - Tests requiring credentials should load them from a `test_config.yml`.
  - All tests should be run under **testnet environments** to avoid impacting real accounts.
- Ensure tests are:
  - **Idempotent** and can be re-run without manual cleanup.
  - **Isolated** from live trading systems.
- Naming:
  - Use descriptive test names indicating the API and scenario being tested.
  - Place tests in `*_test.go` in the corresponding package.

---

## Conventions and Best Practices

- **Use `context.Context` in all public APIs for cancellation and timeout control.**
- **Consistent logging using `log/slog`** for structured, level-based logs.
- Avoid global variables (except for constants).
- Exported functions, methods, and types **must have clear documentation**.
- Respect Go idioms for error handling and simplicity.
- Structure WebSocket reconnections with backoff and error handling.
- Use environment variables for secrets (`API_KEY`, `API_SECRET`) and avoid hardcoding sensitive information.
- Maintain **a clear separation of concerns**:
  - Low-level communication in `ws.go`/`request.go`.
  - Business logic orchestration in `ws_client.go`/`client.go`.

---


## Maintainers

- [Yi Te], Lead Developer

Contact: [coastq22889@icloud.com]

---

## References

- [Go Documentation](https://go.dev/doc/)
- [golangci-lint](https://golangci-lint.run/)
- Exchange-specific API Documentation (Binance, Bybit, OKX, etc.)

---

## Example Commands

### Build & Lint

```bash
go build ./...
golangci-lint run
gofmt -s -w .
```
