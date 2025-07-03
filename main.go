package main

import (
	"encoding/json"
	"fmt"
	"funding-arb-be/utils"
	"log"
	"net/http"
	"strings"

	"funding-arb-be/services"
	"github.com/gorilla/mux"
)

func main() {
	utils.LoadEnv()
	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}).Methods("GET")

	r.HandleFunc("/test-auth/{exchange}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		exchange := vars["exchange"]

		switch exchange {
		case "binance":
			services.TestBinanceAuth()
		case "bybit":
			services.TestBybitAuth()
		case "okx":
			services.TestOkxAuth()
		case "hyperliquid":
			services.TestHyperliquidReadOnly()
		default:
			http.Error(w, "Unknown exchange", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Auth test sent to %s\n", exchange)
	}).Methods("POST")

	r.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
		var tradeReq services.TradeRequest
		err := json.NewDecoder(r.Body).Decode(&tradeReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received trade: %+v\n", tradeReq)

		switch strings.ToLower(tradeReq.Exchange) {
		case "bybit":
			services.ExecuteBybitTrade(tradeReq)
		default:
			http.Error(w, "Unsupported exchange", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "Trade request received")
	}).Methods("POST")

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
