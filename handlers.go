package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"t7r.dev/demo/vrp"
)

func createMandateHandler(client *vrp.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		mandateResp, err := client.CreateMandate(vrp.CreateMandateRequest{
			Mandate: vrp.Mandate{
				Type: vrp.Sweeping,
				ProviderSelection: vrp.MandateProviderSelection{
					Type:       vrp.Preselected, // preselected to avoid implementing user selection panel for now
					ProviderID: "ob-natwest-vrp-sandbox",
				},
				Beneficiary: vrp.Beneficiary{
					Type:              vrp.External,
					AccountHolderName: "Big Bucks LTD",
					AccountIdentifier: vrp.PaymentAccountIdentifier{
						Type:          vrp.SCAN,
						SortCode:      "012345",
						AccountNumber: "12345678",
					},
				},
			},
			Currency: vrp.GBP,
			User: vrp.User{
				ID:    uuid.New().String(),
				Name:  "A N Other",
				Email: "an.other@example.com",
			},
			Constraints: vrp.Constraints{
				ValidFrom:               vrp.ZuluTime{Time: time.Now()},
				ValidTo:                 vrp.ZuluTime{Time: time.Now().Add(24 * time.Hour)},
				MaximumIndividualAmount: 1,
				PeriodicLimits: vrp.PeriodicLimits{
					Day: &vrp.PeriodicLimit{
						MaximumAmount:   5,
						PeriodAlignment: vrp.Calendar,
					},
				},
			},
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "mandate create error %w", err)
			return
		}

		fmt.Println("Created Mandate", mandateResp.ID)

		authResp, err := client.StartAuth(mandateResp.ID, vrp.StartAuthFlowRequest{
			Redirect: vrp.Redirect{
				ReturnURI: fmt.Sprintf("%s/callback", os.Getenv("BASE_HOST")),
			},
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "authorization error %w", err)
			return
		}

		c.Redirect(http.StatusSeeOther, authResp.AuthorizationFlow.Actions.Next.URI)
	}
}

func createPaymentHandler(client *vrp.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		amount, err := strconv.Atoi(c.PostForm("amount"))
		if err != nil {
			c.String(http.StatusBadRequest, "bad amount %w", err)
			return
		}

		resp, err := client.CreatePayment(vrp.CreatePaymentRequest{
			AmountInMinor: amount,
			Currency:      vrp.GBP,
			PaymentMethod: vrp.PaymentMethod{
				Type:      vrp.PaymentMethodMandate,
				MandateID: c.PostForm("mandate_id"),
			},
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "payment create error %w", err)
			return
		}
		fmt.Println("Created Payment", resp.ID)
		c.HTML(http.StatusOK, "payment.html", gin.H{"mandate_id": c.PostForm("mandate_id")})
	}
}
