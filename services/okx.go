package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func TestOkxAuth() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_SECRET_KEY")
	passphrase := os.Getenv("OKX_PASSPHRASE")

	timestamp := time.Now().UTC().Format(time.RFC3339)
	method := "GET"
	requestPath := "/api/v5/account/balance"
	body := ""

	// Sign string: timestamp + method + requestPath + body
	signStr := timestamp + method + requestPath + body

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	url := "https://www.okx.com" + requestPath
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("OK-ACCESS-KEY", apiKey)
	req.Header.Set("OK-ACCESS-SIGN", signature)
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", passphrase)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("OKX Auth Response:", string(bodyBytes))
}

func ExecuteOkxTrade(req TradeRequest) {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_SECRET_KEY")
	passphrase := os.Getenv("OKX_PASSPHRASE")

	side := "buy"
	if strings.ToLower(req.Side) == "short" {
		side = "sell"
	}

	// fetch minSz and lotSz
	minSzStr, lotSzStr, err := GetOkxSymbolSizeConfig(req.Symbol)
	if err != nil {
		fmt.Println("Failed to fetch size config:", err)
		return
	}

	// parse to float
	minSz, _ := strconv.ParseFloat(minSzStr, 64)
	lotSz, _ := strconv.ParseFloat(lotSzStr, 64)

	// adjust amount to valid multiple
	roundedQty := math.Floor(req.Amount/lotSz) * lotSz
	if roundedQty < minSz {
		fmt.Printf("Adjusted quantity %.8f is below minSz %.8f\n", roundedQty, minSz)
		return
	}

	order := map[string]string{
		"instId":  req.Symbol,
		"tdMode":  "cross",
		"side":    side,
		"posSide": req.Side,
		"ordType": "market",
		"sz":      fmt.Sprintf("%.8f", roundedQty),
	}

	body, err := json.Marshal(order)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	method := "POST"
	requestPath := "/api/v5/trade/order"

	signStr := timestamp + method + requestPath + string(body)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	url := "https://www.okx.com" + requestPath
	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("OK-ACCESS-KEY", apiKey)
	httpReq.Header.Set("OK-ACCESS-SIGN", signature)
	httpReq.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	httpReq.Header.Set("OK-ACCESS-PASSPHRASE", passphrase)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Error placing order:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("OKX trade response:", string(bodyBytes))
}

func GetOkxSymbolSizeConfig(symbol string) (minSz string, lotSz string, err error) {
	url := "https://www.okx.com/api/v5/public/instruments?instType=SWAP"
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch OKX instruments: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Code string          `json:"code"`
		Data []OkxInstrument `json:"data"`
	}

	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse OKX response: %v", err)
	}

	for _, inst := range parsed.Data {
		if inst.InstId == symbol {
			return inst.MinSz, inst.LotSz, nil
		}
	}

	return "", "", fmt.Errorf("symbol %s not found", symbol)
}
