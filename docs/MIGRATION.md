# Fabric → Fabric-X Migration Analysis

## Overview

This document analyses what would be required to migrate the chaincode in this
simulator from classic Hyperledger Fabric to Fabric-X.  It is the core
intellectual contribution of this project — separate from the code — and
directly addresses the goal of the LFX mentorship project.

---

## What is Fabric-X?

Fabric-X is a rearchitected version of Hyperledger Fabric that introduces:

- An **orderer-centric execution model** — transaction ordering and execution
  are more tightly coupled, removing the peer endorsement round-trip in some
  configurations.
- An **external chaincode launcher** as the primary deployment model — chaincode
  runs as an independent process and communicates with the peer over gRPC,
  rather than being managed inside the peer process.
- A **simplified transaction flow** — aiming to reduce the steps between
  client submission and ledger commit.

---

## What stays the same (portable layer)

The chaincode business logic in `chaincode/asset_contract.go` is **fully
portable** to Fabric-X without modification.  Specifically:

| Component | Why it's portable |
|---|---|
| `PutState` / `GetState` / `DelState` | Core world-state API is preserved in Fabric-X |
| Asset serialisation (JSON) | Data model is independent of transport |
| Business logic (CreateAsset, TransferAsset, etc.) | Pure Go functions with no peer dependency |
| Error handling patterns | Unchanged between Fabric and Fabric-X |

**Key insight**: The `StubInterface` in this simulator defines exactly the
portable surface.  Any chaincode that only depends on `PutState`, `GetState`,
`DelState`, and `GetTxID` will migrate to Fabric-X with zero business-logic
changes.

---

## What changes (migration surface)

### 1. Shim bootstrapping

**Classic Fabric**
```go
// main.go in classic Fabric chaincode
func main() {
    cc := new(AssetContract)
    if err := contractapi.NewChaincode(cc); err != nil {
        log.Panicf("Error creating chaincode: %v", err)
    }
}
```
The peer shim (`fabric-chaincode-go/shim`) manages the gRPC connection to the
peer internally.  The developer calls `shim.Start()` and the framework handles
everything else.

**Fabric-X**
Fabric-X uses the *external chaincode launcher* model as default.  The
chaincode must:
1. Connect to the peer's chaincode gRPC endpoint explicitly using connection
   details from an injected `CHAINCODE_SERVER_ADDRESS` and `CHAINCODE_ID`
   environment variable.
2. Implement the `ChaincodeServer` interface and serve requests over that
   connection.

Migration action: replace `shim.Start(cc)` with a gRPC server setup:
```go
// Fabric-X style bootstrap (pseudo-code)
server := &shim.ChaincodeServer{
    CCID:    os.Getenv("CHAINCODE_ID"),
    Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
    CC:      new(AssetContract),
    TLSProps: shim.TLSProperties{Disabled: true}, // TLS config
}
server.Start()
```

### 2. Deployment lifecycle

**Classic Fabric**

```bash
peer lifecycle chaincode package mycc.tar.gz ...
peer lifecycle chaincode install mycc.tar.gz
peer lifecycle chaincode approveformyorg ...
peer lifecycle chaincode commit ...
```

**Fabric-X**
The chaincode is deployed as an independent service (Docker container, k8s
pod, or bare process).  The peer references it by its `CCID` and address.
There is no `install` step that pushes a binary to the peer.

Migration action: write a `Dockerfile` for the chaincode, deploy it
independently, and register its address with the Fabric-X peer configuration.

### 3. Endorsement policy and execution model

Classic Fabric requires multiple peers to *endorse* (simulate and sign) a
transaction before it is sent to the orderer.  Fabric-X's streamlined model
changes when and how endorsement happens.

Migration action: review endorsement policies and update channel configuration
for Fabric-X's transaction flow.  Chaincode code itself is unaffected.

---

## Migration effort estimate

| Area | Effort | Code change required |
|---|---|---|
| Business logic (contracts) | None | No |
| State model (JSON assets) | None | No |
| Shim bootstrap (main.go) | Low | ~20 lines |
| Deployment scripts | Medium | New Dockerfile + peer config |
| Endorsement policy config | Medium | Channel config update |
| Test suite | Low | Update connection setup |

---

## What this simulator demonstrates

This project proves that the **portable layer is real and testable**:

- `chaincode/asset_contract.go` runs against `MockStub` (this simulator) and
  would run identically against the real Fabric peer shim or a Fabric-X shim,
  because all three satisfy `StubInterface`.
- The `MockStub` itself is a concrete example of what a minimal Fabric-X
  compatibility shim would need to implement.
- The transaction history feature mirrors `GetHistoryForKey`, a core Fabric API
  whose continued availability in Fabric-X is a key migration question.

---

## Recommendations for the LFX PoC

Based on this analysis, the following areas are recommended for exploration
during the mentorship:

1. **Validate the portable layer against fabric-samples/test-network**
   Run `asset-transfer-basic` against a live Fabric network to confirm
   that chaincode depending only on PutState/GetState/DelState requires
   zero modification for the FSC port.

2. **Map the bootstrap diff explicitly**
   Document the exact line-by-line difference between classic Fabric
   `shim.Start()` and the Fabric-X external launcher gRPC server setup.
   This is the clearest migration surface identified in this analysis.

3. **Catalogue stub methods beyond the core four**
   Real-world chaincode (CC-Tools, Fabric Private Chaincode) uses rich
   queries, private data collections, and events. Each of these needs
   a Fabric-X equivalent identified or a migration workaround documented.

4. **Produce a category-based migration guide**
   Group stub methods by migration effort: portable as-is, requires
   bootstrap change only, requires architectural redesign. This would
   serve as a practical reference for teams migrating existing Fabric
   applications.

---

## References

- [Hyperledger Fabric External Chaincode Launcher](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html)
- [fabric-contract-api-go](https://github.com/hyperledger/fabric-contract-api-go)
- [Fabric-X repository](https://github.com/hyperledger/fabric-x)
- [CC-Tools framework](https://github.com/hyperledger-labs/cc-tools)
