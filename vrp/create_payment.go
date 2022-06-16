package vrp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) CreatePayment(req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	httpResp, err := c.DoAPIRequest("POST", "payments", req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error body %w", err)
	}

	if httpResp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error %s %s", httpResp.Status, string(respBody))
	}

	var resp CreatePaymentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode success body %w %s", err, string(respBody))
	}
	return &resp, nil
}

type CreatePaymentResponse struct {
	ID   string `json:"id"`
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	ResourceToken string `json:"resource_token"`
}

type CreatePaymentRequest struct {
	AmountInMinor int           `json:"amount_in_minor"`
	Currency      Currency      `json:"currency"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type PaymentMethod struct {
	Type      PaymentMethodType `json:"type"`
	MandateID string            `json:"mandate_id"`
}
type PaymentMethodType = string

const (
	PaymentMethodMandate PaymentMethodType = "mandate"
)
