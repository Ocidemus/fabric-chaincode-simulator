# Fabric Chaincode Execution Simulator

A local execution harness for Hyperledger Fabric-style Go chaincode — built to
explore the migration path from classic Fabric to **Fabric-X**.

---

## Why this project exists

The LFX Fabric-X mentorship project asks: *can existing Fabric chaincode run on
Fabric-X, and what would migration require?*

This simulator answers that question from the bottom up:

1. It implements the **same interface** (`StubInterface`) that Fabric's peer
   shim exposes to chaincode at runtime.
2. It runs real chaincode business logic against that interface locally, with
   no Fabric network required.
3. The `docs/MIGRATION.md` document analyses exactly which parts of the
   chaincode are portable and which parts need to change for Fabric-X.

---

## Architecture

```
Client (Postman / curl)
         │
         ▼
  POST /invoke  ──────────────────────────────────────────────────────────┐
         │                                                                │
         ▼                                                                │
  InvokeHandler                                                           │
  (handlers/invoke.go)                                                    │
         │  sets TxID, dispatches function                               │
         ▼                                                                │
  AssetContract                                                           │
  (chaincode/asset_contract.go)  ◄── PORTABLE: identical to real Fabric  │
         │  calls PutState / GetState / DelState                         │
         ▼                                                                │
  MockStub                                                                │
  (stub/mock_stub.go)  ◄── implements same interface as Fabric peer shim │
         │  in-memory world state + append-only transaction history      │
         └────────────────────────────────────────────────────────────────┘

GET /history/{id}    →  returns per-asset transaction log (simulates GetHistoryForKey)
GET /transactions    →  returns full ledger log
```

---

## Project structure

```
fabric-chaincode-simulator/
├── main.go                      HTTP server, route registration
├── chaincode/
│   └── asset_contract.go        Chaincode business logic (portable to Fabric-X)
├── stub/
│   └── mock_stub.go             MockStub — implements ChaincodeStubInterface subset
├── handlers/
│   ├── invoke.go                POST /invoke — function dispatch
│   └── history.go               GET /history, GET /transactions
├── models/
│   ├── asset.go                 Asset struct (on-chain data model)
│   └── transaction.go           Transaction record (ledger history entry)
├── utils/
│   └── json.go                  HTTP response helpers
└── docs/
    └── MIGRATION.md             Fabric → Fabric-X migration analysis
```

---

## Running locally

```bash
git clone <repo>
cd fabric-chaincode-simulator
go run main.go
# Server starts on :8080
```

---

## API reference

### POST /invoke

Execute a chaincode function.

```json
{
  "function": "CreateAsset",
  "args": {
    "id": "asset1",
    "owner": "Alice",
    "value": "500"
  }
}
```

Supported functions:

| Function | Required args | Description |
|---|---|---|
| `CreateAsset` | `id`, `owner`, `value` | Create a new asset |
| `TransferAsset` | `id`, `newOwner` | Change asset owner |
| `UpdateValue` | `id`, `value` | Update asset value |
| `DeleteAsset` | `id` | Remove asset from world state |
| `GetAsset` | `id` | Read asset data |
| `AssetExists` | `id` | Check if asset exists |

### GET /history/{assetID}

Returns the full transaction history for an asset, simulating Fabric's
`GetHistoryForKey` stub method.

### GET /transactions

Returns all transactions across all assets — the full simulated ledger log.

---

## Key design decisions

### StubInterface is the migration boundary

`chaincode/asset_contract.go` depends only on `StubInterface`:

```go
type StubInterface interface {
    PutState(key string, value []byte) error
    GetState(key string) ([]byte, error)
    DelState(key string) error
    GetTxID() string
}
```

Both `MockStub` (this simulator) and a real Fabric peer stub implement this
interface.  The contract code does not import the Fabric SDK at all — which
means it is testable, portable, and ready for Fabric-X with only the bootstrap
`main.go` needing to change.

### Append-only transaction history

Every invocation — success or failure — is recorded to an immutable log inside
`MockStub`.  This mirrors the Fabric ledger's behaviour where committed
transactions are permanently recorded and queryable via `GetHistoryForKey`.

---

## Connection to Fabric-X

See `docs/MIGRATION.md` for the full analysis.  Short version:

| Layer | Portable? | Migration effort |
|---|---|---|
| Business logic | ✅ Yes | None |
| State model | ✅ Yes | None |
| Shim bootstrap | ❌ Changes | ~20 lines (gRPC server setup) |
| Deployment | ❌ Changes | Dockerfile + peer config |

---

## Relation to real Fabric chaincode

In a production Fabric deployment, `AssetContract` would be used like this:

```go
// production main.go (classic Fabric)
import "github.com/hyperledger/fabric-contract-api-go/contractapi"

func main() {
    cc, _ := contractapi.NewChaincode(&chaincode.AssetContract{})
    cc.Start()
}
```

The `contractapi` framework injects the real `ChaincodeStubInterface` — which
is a superset of the `StubInterface` defined here.  The business logic methods
do not change.

---

## Skills demonstrated

- Intermediate Go: interfaces, JSON marshalling, struct embedding, concurrency-safe maps
- Backend architecture: separation of concerns, dependency injection via interface
- Blockchain fundamentals: world state, transaction lifecycle, endorsement model
- Fabric-specific knowledge: chaincode contract API patterns, stub interface, ledger history
- Migration thinking: identifying the portable vs. transport-specific layers
