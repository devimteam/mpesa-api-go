package test

import (
	"bytes"
	"encoding/base64"
	"flag"
	"testing"
	"time"

	"github.com/devimteam/mpesa-api-go"
	"gotest.tools/assert"
)

var (
	flagKey    = flag.String("key", "", "Daraja MPESA Consumer Key")
	flagSecret = flag.String("secret", "", "Daraja MPESA Consumer Secret")
)

var testMPESAService *mpesa.Service

func TestMain(t *testing.M) {
	flag.Parse()
	testMPESAService = mpesa.New(*flagKey, *flagSecret, mpesa.SandboxEndpoint)
	t.Run()
}

const (
	shortCode1            = "600351"
	initiatorName         = "testapi0351"
	securityCredential    = "Safaricom351!"
	shortCode2            = "600000"
	testMSISDN            = "254708374149"
	mpesaOnlineShortcode  = "174379"
	mpesaOnlinePasskey    = "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
	initiatorSecurityCred = "EA88UFQE6V9G2HkmuubuLpj/CdySZm1q9YtrXFV0KXw7gSNQONvWR9YhrQin595gpid92PvwatjpOopq+L6sEEi4HdtrfYBEgvW+HUgKMhSrJonl29nunu/t6NbMOiuvUFZ5NYxo1vsLOnAK4useJCTZPFCHCP8TTAE8SLWjOS6fsZcNkhMMU6YHUJq4ptKYWDvW/+EkfHYM/SEJPnZtPhZ6XhnGsQOBTHfpn+XmLz6PAK7L2Y1FvEKJP62Jo7+JIdtxNXVIyV1OVg4kwnoDZv9kF/GdW0tjBRwfaow6VPMuh2e7SKYJT36dm5PznbX/1liI/6z08fcwJrNUV0RMnA=="

	callbackUrl = "http://google.com:80/somepath"
)

func TestService_GenerateNewAccessToken(t *testing.T) {
	token, err := testMPESAService.GenerateNewAccessToken()
	assert.NilError(t, err)
	assert.Assert(t, token != "", "token is empty")
}

func TestService_B2CRequest(t *testing.T) {
	paymentResp, err := testMPESAService.B2CRequest(mpesa.B2C{
		InitiatorName:      initiatorName,
		SecurityCredential: initiatorSecurityCred,
		CommandID:          "PromotionPayment",
		Amount:             "10",
		PartyA:             shortCode1,
		PartyB:             testMSISDN,
		Remarks:            "auto-testing",
		QueueTimeOutURL:    callbackUrl,
		ResultURL:          callbackUrl,
		Occasion:           "",
	})
	assert.NilError(t, err)
	assert.Assert(t, paymentResp != nil, "response is nil")
	assert.Assert(t, paymentResp.ConversationID != "", "ConversationID is empty")
	assert.Assert(t, paymentResp.OriginatorConversationID != "", "OriginatorConversationID is empty")
	assert.Assert(t, paymentResp.ResponseCode != "", "ResponseCode is empty")
	assert.Assert(t, paymentResp.ResponseCode == "0", "ResponseCode is not zero")
	assert.Assert(t, paymentResp.ResponseDescription != "", "ResponseDescription is empty")
}

func TestService_MPESAOnlinePayment(t *testing.T) {
	now := mpesa.Timestamp(time.Now())
	var password bytes.Buffer
	base64.NewEncoder(base64.StdEncoding, &password).Write([]byte(mpesaOnlineShortcode + mpesaOnlinePasskey + now))
	paymentResp, err := testMPESAService.MPESAOnlinePayment(mpesa.Payment{
		BusinessShortCode: mpesaOnlineShortcode,
		Password:          password.String(),
		Timestamp:         mpesa.Timestamp(time.Now()),
		TransactionType:   "CustomerPayBillOnline",
		Amount:            "1",
		PartyA:            testMSISDN,
		PartyB:            mpesaOnlineShortcode,
		PhoneNumber:       testMSISDN,
		CallBackURL:       callbackUrl,
		AccountReference:  "auto-testing", // should not be empty
		TransactionDesc:   "auto-testing", // should not be empty
	})
	assert.NilError(t, err)
	assert.Assert(t, paymentResp != nil, "response is nil")
	assert.Assert(t, paymentResp.MerchantRequestID != "", "MerchantRequestID is empty")
	assert.Assert(t, paymentResp.CheckoutRequestID != "", "CheckoutRequestID is empty")
	assert.Assert(t, paymentResp.ResponseCode != "", "ResponseCode is empty")
	assert.Assert(t, paymentResp.ResponseCode == "0", "ResponseCode is not zero")
	assert.Assert(t, paymentResp.ResponseDescription != "", "ResponseDescription is empty")
	assert.Assert(t, paymentResp.CustomerMessage != "", "CustomerMessage is empty")
}

func TestService_C2BSimulation(t *testing.T) {
	paymentResp, err := testMPESAService.C2BSimulation(mpesa.C2B{
		ShortCode:     shortCode1,
		CommandID:     "CustomerPayBillOnline",
		Amount:        "10",
		Msisdn:        testMSISDN,
		BillRefNumber: "auto-testing",
	})
	assert.NilError(t, err)
	assert.Assert(t, paymentResp != nil, "response is nil")
	assert.Assert(t, paymentResp.OriginatorConversationID == "", "OriginatorConversationID is empty")
	assert.Assert(t, paymentResp.ConversationID != "", "ConversationID is empty")
	assert.Assert(t, paymentResp.ResponseCode == "", "ResponseCode is empty")
	assert.Assert(t, paymentResp.ResponseDescription != "", "ResponseDescription is empty")
}

func TestService_C2BRegisterURL(t *testing.T) {
	paymentResp, err := testMPESAService.C2BRegisterURL(mpesa.C2BRegisterURL{
		ShortCode:       shortCode1,
		ResponseType:    "Cancelled",
		ConfirmationURL: callbackUrl,
		ValidationURL:   callbackUrl,
	})
	assert.NilError(t, err)
	assert.Assert(t, paymentResp != nil, "response is nil")
	assert.Assert(t, paymentResp.OriginatorConversationID == "", "OriginatorConversationID is empty")
	assert.Assert(t, paymentResp.ConversationID == "", "ConversationID is empty")
	assert.Assert(t, paymentResp.ResponseCode == "", "ResponseCode is empty")
	assert.Assert(t, paymentResp.ResponseDescription != "", "ResponseDescription is empty")
	assert.Assert(t, paymentResp.ResponseDescription == "success", "ResponseDescription is not 'success'")
}
