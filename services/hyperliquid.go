package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func TestHyperliquidReadOnly() {
	payload := map[string]interface{}{
		"type": "portfolio",
		"user": "0x9389870e86daa280ba8211e4c95303f4ca435298",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.hyperliquid.xyz/info", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Request creation error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response:", string(respBody))
}
