package vrp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

func (c *Client) StartAuth(id string, req StartAuthFlowRequest) (*StartAuthFlowResponse, error) {
	endpoint := path.Join("mandates", id, "authorization-flow")
	httpResp, err := c.DoAPIRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error body %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error %s %s", httpResp.Status, string(respBody))
	}

	var resp StartAuthFlowResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode success body %w %s", err, string(respBody))
	}
	return &resp, nil
}

type StartAuthFlowResponse struct {
	Status            string            `json:"status"`
	AuthorizationFlow AuthorizationFlow `json:"authorization_flow"`
}

type AuthorizationFlow struct {
	Actions Actions `json:"actions"`
}

type Actions struct {
	Next Next `json:"next"`
}

type Next struct {
	URI string `json:"uri"`
}

type StartAuthFlowRequest struct {
	ProviderSelection ProviderSelection `json:"provider_selection"`
	Redirect          Redirect          `json:"redirect"`
}

type ProviderSelection struct{}

type Redirect struct {
	ReturnURI string `json:"return_uri"`
}
