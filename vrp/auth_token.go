package vrp

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

func (c *Client) GetHttpClient() *http.Client {
	if c.Client == nil {
		config := clientcredentials.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			TokenURL:     fmt.Sprintf("https://auth.%s/connect/token", c.Host),
			Scopes:       []string{"payments", "recurring_payments:sweeping"},
		}
		c.Client = config.Client(context.Background())
	}
	return c.Client
}
