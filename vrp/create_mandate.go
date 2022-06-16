package vrp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *Client) CreateMandate(req CreateMandateRequest) (*CreateMandateResponse, error) {
	httpResp, err := c.DoAPIRequest("POST", "mandates", req)
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

	var resp CreateMandateResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode success body %w %s", err, string(respBody))
	}
	return &resp, nil
}

type CreateMandateResponse struct {
	ID   string `json:"id"`
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	ResourceToken string `json:"resource_token"`
}

type CreateMandateRequest struct {
	Mandate     Mandate     `json:"mandate"`
	Currency    Currency    `json:"currency"`
	User        User        `json:"user"`
	Constraints Constraints `json:"constraints"`
}

type Mandate struct {
	Type              MandateType              `json:"type"`
	ProviderSelection MandateProviderSelection `json:"provider_selection"`
	Beneficiary       Beneficiary              `json:"beneficiary"`
}

type MandateProviderSelection struct {
	Type MandateProviderSelectionType `json:"type"`

	Filter *ProviderFilter `json:"filter,omitempty"` // only on UserSelected flow (required)

	ProviderID string    `json:"provider_id,omitempty"` // only on Preselected flow (required)
	Remitter   *Remitter `json:"remitter,omitempty"`    // only on Preselected flow (optional)
}

type ProviderFilter struct {
	Countries        []string `json:"countries,omitempty"`
	ReleaseChannel   string   `json:"release_channel,omitempty"`
	CustomerSegments []string `json:"customer_segments,omitempty"`
	ProviderIDs      []string `json:"provider_ids,omitempty"`
	Excludes         struct {
		ProviderIDs []string `json:"provider_ids,omitempty"`
	} `json:"excludes,omitempty"`
}

type Remitter struct {
	AccountHolderName string                   `json:"account_holder_name"`
	AccountIdentifier PaymentAccountIdentifier `json:"account_identifier"`
}

type PaymentAccountIdentifier struct {
	Type PaymentAccountIdentifierType `json:"type"`

	IBAN string `json:"iban,omitempty"` // only on IBAN flow (required)

	SortCode      string `json:"sort_code,omitempty"`      // only on SCAN flow (required)
	AccountNumber string `json:"account_number,omitempty"` // only on SCAN flow (required)
}

type PaymentAccountIdentifierType = string

const (
	SCAN PaymentAccountIdentifierType = "sort_code_account_number"
	IBAN PaymentAccountIdentifierType = "iban"
)

type MandateType = string

const (
	Sweeping   MandateType = "sweeping"
	Commercial MandateType = "commercial"
)

type MandateProviderSelectionType = string

const (
	UserSelected MandateProviderSelectionType = "user_selected"
	Preselected  MandateProviderSelectionType = "preselected"
)

type Beneficiary struct {
	Type              BeneficiaryType          `json:"type"`
	AccountHolderName string                   `json:"account_holder_name,omitempty"`
	MerchantAccountID string                   `json:"merchant_account_id,omitempty"`
	AccountIdentifier PaymentAccountIdentifier `json:"account_identifier,omitempty"`
}

type BeneficiaryType = string

const (
	External BeneficiaryType = "external_account"
	Merchant BeneficiaryType = "merchant_account"
)

type Currency = string

const (
	GBP Currency = "GBP"
	EUR Currency = "EUR"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type Constraints struct {
	ValidFrom               ZuluTime       `json:"valid_from"`
	ValidTo                 ZuluTime       `json:"valid_to"`
	MaximumIndividualAmount int            `json:"maximum_individual_amount"`
	PeriodicLimits          PeriodicLimits `json:"periodic_limits"`
}

type PeriodicLimits struct {
	Day       *PeriodicLimit `json:"day,omitempty"`
	Week      *PeriodicLimit `json:"week,omitempty"`
	Fortnight *PeriodicLimit `json:"fortnight,omitempty"`
	Month     *PeriodicLimit `json:"month,omitempty"`
	HalfYear  *PeriodicLimit `json:"half_year,omitempty"`
	Year      *PeriodicLimit `json:"year,omitempty"`
}

type PeriodicLimit struct {
	MaximumAmount   int             `json:"maximum_amount"`
	PeriodAlignment PeriodAlignment `json:"period_alignment"`
}
type PeriodAlignment = string

const (
	Calendar PeriodAlignment = "calendar"
	Consent  PeriodAlignment = "consent"
)

type ZuluTime struct {
	Time time.Time
}

func (z ZuluTime) MarshalText() ([]byte, error) {
	return []byte(z.Time.UTC().Format("2006-01-02T15:04:05.000Z")), nil
}
