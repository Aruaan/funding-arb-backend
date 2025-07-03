package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func TestBybitAuth() {
	apiKey := os.Getenv("BYBIT_API_KEY")
	secretKey := os.Getenv("BYBIT_SECRET_KEY")
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	recvWindow := "5000"

	params := "accountType=UNIFIED"
	signingString := timestamp + apiKey + recvWindow + params

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signingString))
	signature := hex.EncodeToString(h.Sum(nil))

	url := fmt.Sprintf("https://api.bybit.com/v5/account/wallet-balance?%s", params)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-BAPI-API-KEY", apiKey)
	req.Header.Set("X-BAPI-TIMESTAMP", timestamp)
	req.Header.Set("X-BAPI-RECV-WINDOW", recvWindow)
	req.Header.Set("X-BAPI-SIGN", signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Bybit Auth Response:", string(body))
}

func ExecuteBybitTrade(req TradeRequest) {
	fmt.Printf("Placing bybit %s on %s for %.2f\n", req.Side, req.Symbol, float64(req.Amount))
}
