package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
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
