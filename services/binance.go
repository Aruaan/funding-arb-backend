package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func TestBinanceAuth() {
	baseURL := "https://fapi.binance.com"
	endpoint := "/fapi/v3/account"
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")
	timestamp := time.Now().UnixMilli()
	params := fmt.Sprintf("timestamp=%d", timestamp)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(params))
	signature := hex.EncodeToString(h.Sum(nil))

	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", baseURL, endpoint, params, signature)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Binance Futures Auth Response:", string(body))
}

func ExecuteBinanceTrade(req TradeRequest) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	side := "BUY"
	if strings.ToLower(req.Side) == "short" {
		side = "SELL"
	}

	timestamp := time.Now().UnixMilli()

	params := url.Values{}
	params.Set("symbol", req.Symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("quantity", fmt.Sprintf("%.3f", req.Amount))
	params.Set("timestamp", fmt.Sprintf("%d", timestamp))

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(params.Encode()))
	signature := hex.EncodeToString(h.Sum(nil))
	params.Set("signature", signature)

	urlStr := fmt.Sprintf("https://fapi.binance.com/fapi/v1/order?%s", params.Encode())
	httpReq, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	httpReq.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Error placing order:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Binance trade response:", string(body))
}
