package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	tlsigning "github.com/Truelayer/truelayer-signing/go"
	"github.com/gin-gonic/gin"
)

func VerifiedSignature(c *gin.Context) {
	body, err := verifyWebhook(c.Request)
	if err != nil {
		fmt.Println("Bad webhook recieved")
		c.Status(http.StatusUnauthorized)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	c.Next()
}

func verifyWebhook(r *http.Request) ([]byte, error) {
	tlSignature := r.Header.Get("Tl-Signature")
	if len(tlSignature) == 0 {
		return nil, fmt.Errorf("missing Tl-Signature header")
	}

	jwsHeader, err := tlsigning.ExtractJwsHeader(tlSignature)
	if err != nil {
		return nil, err
	}
	if len(jwsHeader.Jku) == 0 {
		return nil, fmt.Errorf("jku missing")
	}

	defer r.Body.Close()
	webhookBody, err := io.ReadAll(r.Body)
	if err != nil {
		return webhookBody, fmt.Errorf("webhook body missing")
	}

	// ensure jku is an expected TrueLayer url
	if jwsHeader.Jku != fmt.Sprintf("https://webhooks.%s/.well-known/jwks", getHost()) {
		return webhookBody, fmt.Errorf("unpermitted jku %s", jwsHeader.Jku)
	}

	// fetch jwks (should be cached according to cache-control headers)
	resp, err := http.Get(jwsHeader.Jku)
	if err != nil {
		return webhookBody, fmt.Errorf("failed to fetch jwks")
	}
	defer resp.Body.Close()
	jwks, err := io.ReadAll(resp.Body)
	if err != nil {
		return webhookBody, fmt.Errorf("jwks missing")
	}

	// verify signature using the jwks
	return webhookBody, tlsigning.
		VerifyWithJwks(jwks).
		Method(http.MethodPost).
		Path(r.RequestURI).
		Headers(getHeadersMap(r.Header)).
		Body(webhookBody).
		Verify(tlSignature)
}

func getHeadersMap(requestHeaders map[string][]string) map[string][]byte {
	headers := make(map[string][]byte)
	for key, values := range requestHeaders {
		// take first value
		headers[key] = []byte(values[0])
	}
	return headers
}
