
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ocidemus/fabric-chaincode-simulator/chaincode"
	"github.com/Ocidemus/fabric-chaincode-simulator/handlers"
	"github.com/Ocidemus/fabric-chaincode-simulator/stub"
)

func main() {
	mockStub := stub.NewMockStub()

	contract := &chaincode.AssetContract{}

	mux := http.NewServeMux()

	mux.Handle("/invoke", handlers.NewInvokeHandler(mockStub, contract))

	mux.Handle("/history/", handlers.NewHistoryHandler(mockStub))

	mux.Handle("/transactions", handlers.NewAllTxHandler(mockStub))

	port := ":8080"
	fmt.Printf("Fabric Chaincode Simulator running on %s\n", port)
	fmt.Println("POST /invoke        — execute a chaincode function")
	fmt.Println("GET  /history/{id}  — asset transaction history")
	fmt.Println("GET  /transactions  — full ledger log")

	log.Fatal(http.ListenAndServe(port, mux))
}
