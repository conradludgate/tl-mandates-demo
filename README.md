# tl-mandates-demo

Demo App written in Go to create mandates and payments using the TrueLayer Payments V3 API.

## Create TL Account

[Create a TL account](https://console.truelayer.com/), create your application and take a note of the client secret.
Next, go to the [Payments](https://console.truelayer.com/payments/v3) settings.
Create your private key and upload the corresponding public key. Keep the private key on hand.

## Ngrok

[Setup a free ngrok account](https://dashboard.ngrok.com/) and download the ngrok program.

Run `ngrox http 8080` and keep it running in the terminal window. Copy the 'Forwarding' URL.

In the TL Console, add `${FORWARNING_URL}/callback` to your application's Settings > Allowed redirect URIs section.
Add `${FORWARDING_URL}/webhook` to your applcation's Payments > Settings > Webhook URI section.

## Run

Save your private key locally to a file named `ec512-private.pem`. Create the following environment variables

```
export TL_KID="{{ YOUR KEY ID }}"
export TL_CLIENT_ID="{{ YOUR CLIENT ID }}"
export TL_CLIENT_SECRET="{{ YOUR CLIENT SECRET }}"
export BASE_HOST="{{ YOUR FORWARDING_URL }}"
export TL_ENVIRONMENT="sandbox"
```

then run

```sh
go run *.go
```

to launch the [gin](https://github.com/gin-gonic/gin) web app.

Navigate to your forwarding url in your browser and you should be able to click the button to create a mandate using the
NatWest Sandbox provider. Enter `123456789012` as the customer account number, and `572436` on the password/pin page. Select
a bank account and click 'Confirm VRP' to finish authorizing the mandate.

One you have the mandate, you can create the payment. The mandate will have a limit of 1p per payment and 5p per day.
Try making six payments. While you're doing this, you should recieve webhook notifications. Check your terminal for logs.
On the sixth payment attempt, you should see a `payment_failure` event.
