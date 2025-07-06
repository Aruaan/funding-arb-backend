package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	apiKey := os.Getenv("BYBIT_API_KEY")
	secretKey := os.Getenv("BYBIT_SECRET_KEY")
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	recvWindow := "5000"

	side := "Buy"
	if strings.ToLower(req.Side) == "short" {
		side = "Sell"
	}

	bodyMap := map[string]interface{}{
		"category":  "linear",
		"symbol":    req.Symbol,
		"side":      side,
		"orderType": "Market",
		"qty":       fmt.Sprintf("%.3f", req.Amount),
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		fmt.Println("Error marshalling body:", err)
		return
	}
	bodyStr := string(bodyBytes)

	// Correct signing string for POST
	signingStr := timestamp + apiKey + recvWindow + bodyStr
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signingStr))
	signature := hex.EncodeToString(h.Sum(nil))

	// Build request
	urlStr := "https://api.bybit.com/v5/order/create"
	httpReq, err := http.NewRequest("POST", urlStr, strings.NewReader(bodyStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-BAPI-API-KEY", apiKey)
	httpReq.Header.Set("X-BAPI-TIMESTAMP", timestamp)
	httpReq.Header.Set("X-BAPI-RECV-WINDOW", recvWindow)
	httpReq.Header.Set("X-BAPI-SIGN", signature)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Error placing order:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Bybit trade response:", string(body))
}
