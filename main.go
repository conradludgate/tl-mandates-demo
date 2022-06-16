package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"t7r.dev/demo/vrp"
)

func main() {
	client := vrp.Client{
		PrivateKey:   getPrivateKey(),
		KID:          os.Getenv("TL_KID"),
		ClientID:     os.Getenv("TL_CLIENT_ID"),
		ClientSecret: os.Getenv("TL_CLIENT_SECRET"),
		Host:         getHost(),
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/callback", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/mandate/%s", c.Query("mandate_id")))
	})
	r.GET("/mandate/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mandate.html", gin.H{"mandate_id": c.Param("id")})
	})
	r.POST("/payment", createPaymentHandler(&client))
	r.POST("/mandate", createMandateHandler(&client))
	r.Use(VerifiedSignature).POST("/webhook", func(c *gin.Context) {
		defer c.Request.Body.Close()
		body, _ := io.ReadAll(c.Request.Body)
		fmt.Println("Recieved", string(body))
	})
	r.Run()
}

func getPrivateKey() []byte {
	privateKey, err := ioutil.ReadFile("ec512-private.pem")
	if err != nil {
		panic(err)
	}
	return privateKey
}

func getHost() string {
	switch os.Getenv("TL_ENVIRONMENT") {
	case "development":
		return "t7r.dev"
	case "sandbox":
		return "truelayer-sandbox.com"
	case "production":
		return "truelayer.com"
	default:
		panic("unknown environment")
	}
}
