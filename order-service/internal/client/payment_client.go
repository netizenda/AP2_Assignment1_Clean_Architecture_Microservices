package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PaymentClient struct {
	client  *http.Client
	baseURL string
}

func NewPaymentClient(baseURL string) *PaymentClient {
	fmt.Printf("[DEBUG] PaymentClient created with baseURL: %s\n", baseURL)
	return &PaymentClient{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: baseURL,
	}
}

func (c *PaymentClient) Authorize(ctx context.Context, orderID string, amount int64) (string, string, error) {
	url := c.baseURL + "/payments"
	fmt.Printf("[DEBUG] Calling Payment Service: %s | orderID=%s | amount=%d\n", url, orderID, amount)

	payload := map[string]any{"order_id": orderID, "amount": amount}
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("[DEBUG] JSON Marshal error: %v\n", err)
		return "", "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("[DEBUG] NewRequest error: %v\n", err)
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("[DEBUG] HTTP Do error: %v\n", err)
		return "", "", err
	}
	defer resp.Body.Close()

	fmt.Printf("[DEBUG] Payment responded with status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[DEBUG] Non-200 status: %d\n", resp.StatusCode)
		return "", "", fmt.Errorf("payment service returned %d", resp.StatusCode)
	}

	var p struct {
		Status        string `json:"status"`
		TransactionID string `json:"transaction_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		fmt.Printf("[DEBUG] JSON Decode error: %v\n", err)
		return "", "", err
	}

	fmt.Printf("[DEBUG] Payment success: status=%s, transaction_id=%s\n", p.Status, p.TransactionID)
	return p.Status, p.TransactionID, nil
}
