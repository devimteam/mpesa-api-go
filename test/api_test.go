package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/devimteam/mpesa-api-go"
)

var testMPESAService = mpesa.New("", "", "")

const (
	shortCode1            = ""
	initiatorName         = ""
	securityCredential    = ""
	shortCode2            = ""
	testMSISDN            = ""
	mpesaOnlineShortcode  = ""
	mpesaOnlinePasskey    = ""
	initiatorSecurityCred = ""

	callbackUrl = ""
)

func TestService_GenerateNewAccessToken(t *testing.T) {
	token, err := testMPESAService.GenerateNewAccessToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}

func TestService_B2CRequest(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go startHttpServer("/", &wg)
	paymentResp, err := testMPESAService.B2CRequest(mpesa.B2C{
		InitiatorName:      initiatorName,
		SecurityCredential: initiatorSecurityCred,
		CommandID:          "BusinessPayment",
		Amount:             "10",
		PartyA:             shortCode1,
		PartyB:             testMSISDN,
		Remarks:            "testing",
		QueueTimeOutURL:    callbackUrl,
		ResultURL:          callbackUrl,
		Occassion:          "",
	})
	wg.Done()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(paymentResp)
	wg.Wait()
}

func TestService_MPESAOnlinePayment(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go startHttpServer("/", &wg)
	paymentResp, err := testMPESAService.MPESAOnlinePayment(mpesa.Payment{
		BusinessShortCode: mpesaOnlineShortcode,
		Password:          mpesaOnlinePasskey,
		Timestamp:         mpesa.Timestamp(time.Now()),
		TransactionType:   "CustomerPayBillOnline",
		Amount:            "100",
		PartyA:            testMSISDN,
		PartyB:            mpesaOnlineShortcode,
		PhoneNumber:       testMSISDN,
		CallBackURL:       callbackUrl,
		AccountReference:  "",
		TransactionDesc:   "testing",
	})
	wg.Done()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(paymentResp)
	wg.Wait()
}

func startHttpServer(prefix string, g *sync.WaitGroup) {
	http.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		defer g.Done()
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Println("received data:", string(data))
		fmt.Fprintf(w, "OK")
	})
	fmt.Println("starting http server on :9999")
	http.ListenAndServe(":9999", nil)
}
