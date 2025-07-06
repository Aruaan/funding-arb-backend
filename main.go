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
		var dummyTrade services.TradeRequest

		switch exchange {
		case "binance":
			dummyTrade = services.TradeRequest{
				Symbol: "BTCUSDT",
				Side:   "long",
				Amount: 0.001, // Binance min lot size for BTC
			}
			services.ExecuteBinanceTrade(dummyTrade)

		case "bybit":
			dummyTrade = services.TradeRequest{
				Symbol: "BTCUSDT",
				Side:   "long",
				Amount: 0.02, // Bybit min size
			}
			services.ExecuteBybitTrade(dummyTrade)

		case "okx":
			dummyTrade = services.TradeRequest{
				Symbol: "BTC-USDT-SWAP",
				Side:   "long",
				Amount: 0.02, // OKX SWAP minSz
			}
			services.ExecuteOkxTrade(dummyTrade)

		default:
			http.Error(w, "Unsupported exchange", http.StatusBadRequest)
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
