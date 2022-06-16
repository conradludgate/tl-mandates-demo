package vrp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	tlsigning "github.com/Truelayer/truelayer-signing/go"
	"github.com/google/uuid"
)

type Client struct {
	PrivateKey   []byte
	KID          string
	ClientID     string
	ClientSecret string
	Host         string
	Client       *http.Client
}

func (c *Client) GetAPIEndpoint(path string) string {
	return fmt.Sprintf("https://api.%s/%s", c.Host, path)
}

func (c *Client) GetAPISignature(method string, path string, idempotency string, body []byte) (string, error) {
	signature, err := tlsigning.SignWithPem(c.KID, c.PrivateKey).
		Method(method).
		Path("/"+path).
		Header("Idempotency-Key", []byte(idempotency)).
		Body(body).
		Sign()
	if err != nil {
		return "", fmt.Errorf("signature %w", err)
	}
	return signature, nil
}

func (c *Client) DoAPIRequest(method string, path string, req interface{}) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal %w", err)
	}
	idempotency := uuid.New().String()

	signature, err := c.GetAPISignature(method, path, idempotency, body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(method, c.GetAPIEndpoint(path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Idempotency-Key", idempotency)
	httpReq.Header.Set("TL-Signature", signature)
	httpReq.Header.Set("Content-Type", "application/json")

	return c.GetHttpClient().Do(httpReq)
}
