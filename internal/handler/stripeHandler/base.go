package stripehandler

import (
	"encoding/json"
	"fmt"
	"ipfast_server/pkg/util/log"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/webhook"
)

type Event = stripe.Event

// dev
var stripeKey string
var endpointSecret string

var (
	CancelURL, SuccessURL string // 支付取消和支付成功的回调地址
)

func Setup() {
	mode := viper.GetString("server.payMode")
	if mode == "" {
		log.Fatalln("⚠️  server.payMode is not set in config file")
	}
	log.Info("支付模式：%s", mode)
	if mode == "dev" {
		stripeKey = "sk_test_51Q8IEuCpfaekutYqxxFUKJPlCWuPEwyAAbt53tZVr0nmnErlpGlNYpMUCYxwUcleQbnbKfqWx6icBDLqygmdPqzk0098dEsrtr"
		endpointSecret = "whsec_d3Ihl0QT5cXotN9z2qGVcZVVzUvD1AsB"
	} else if mode == "prod" {
		stripeKey = "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"
		endpointSecret = "whsec_UU47hkIxc0cYOShs0sqseFbBU61iQ0hV"
	} else {
		log.Fatalln("⚠️  server.payMode is not set correctly in config file")
	}
	serverAddr := viper.GetString("server.serverAddr")
	if serverAddr == "" {
		log.Fatalln("⚠️  server.serverAddr is not set in config file")
	}
	CancelURL = fmt.Sprintf("https://%s/ucenter", serverAddr)
	SuccessURL = fmt.Sprintf("https://%s/ucenter", serverAddr)
}

type PayParam struct {
	Amount        float64 `json:"amount"`
	CommodityName string  `json:"commodity_name"`
	Currency      string  `json:"currency"`
	Quantity      int64   `json:"quantity"`
	Oid           string  `json:"oid"`
}

const (
	CheckoutSessionCompleted             = "checkout.session.completed"
	CheckoutSessionExpired               = "checkout.session.expired"
	CheckoutSessionAsyncPaymentFailed    = "checkout.session.async_payment_failed"
	CheckoutSessionAsyncPaymentSucceeded = "checkout.session.async_payment_succeeded"
)

func Init() {
	stripe.Key = stripeKey
}

func Webhook(payload []byte, stripeSignature string) (eventFlag, oid, pid string, err error) {
	event := stripe.Event{}
	if err = json.Unmarshal(payload, &event); err != nil {
		log.Error("⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		return
	}

	event, err = webhook.ConstructEvent(payload, stripeSignature, endpointSecret)
	if err != nil {
		log.Error("⚠️  Webhook signature verification failed. %v\n", err)
		return
	}
	switch string(event.Type) {
	case CheckoutSessionCompleted:
		// 支付成功
	case CheckoutSessionAsyncPaymentSucceeded:
		// 异步支付成功 类似于信用卡延迟结算这种好像
	case CheckoutSessionAsyncPaymentFailed:
		// 异步支付失败
		err = fmt.Errorf("checkout session async payment failed")
		return
	case CheckoutSessionExpired:
		// 支付订单过期
		err = fmt.Errorf("checkout session expired")
		return
	default:
		// 未处理的事件类型
		err = fmt.Errorf("unhandled event type: %s", event.Type)
		return
	}
	var eventData map[string]interface{}
	if err = json.Unmarshal(event.Data.Raw, &eventData); err != nil {
		err = fmt.Errorf("error parsing event data: %v", err)
		return
	}
	oid, ok := eventData["client_reference_id"].(string)
	if !ok {
		err = fmt.Errorf("client_reference_id not found or not a string")
		return
	}

	pid, ok = eventData["payment_intent"].(string)
	if !ok {
		err = fmt.Errorf("payment_intent not found or not a string")
		return
	}
	eventFlag = string(event.Type)
	return
}

func CreateCheckoutSession(param *PayParam) (string, error) {
	param.Amount = param.Amount * 100 // 转换为美分
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"alipay",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(param.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(param.CommodityName),
					},
					UnitAmount: stripe.Int64(int64(param.Amount)),
				},
				Quantity: stripe.Int64(param.Quantity),
			},
		},
		ClientReferenceID: stripe.String(param.Oid),
		SuccessURL:        stripe.String(SuccessURL + "?oid=" + param.Oid),
		CancelURL:         stripe.String(CancelURL + "?oid=" + param.Oid),
	}

	s, err := session.New(params)

	if err != nil {
		log.Error("CreateCheckoutSession failed: %v", err)
		return "", err
	}
	return s.URL, nil
}
