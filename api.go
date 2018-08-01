package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	SandboxEndpoint    = "https://sandbox.safaricom.co.ke/"
	ProductionEndpoint = "https://api.safaricom.co.ke/"
)

const (
	authHeader        = "Authorization"
	contentTypeHeader = "Content-Type"
	defaultTokenLive  = time.Minute * 45
)

var ErrTokenIsExpired = errors.New("token was expired")

const TimestampLayout = "20060102150405"

func Timestamp(t time.Time) string {
	return t.Format(TimestampLayout)
}

// Service is an Mpesa Service
type Service struct {
	appKey    string
	appSecret string
	endpoint  string

	authHeader string
	// The OAuth access token expires after an hour, after which,
	// you will need to generate another access token.
	// On a production app, use a base64 library of the programming language you are using to build your app to get
	// the Basic Auth string that you will then use to invoke our OAuth API to get an access token.
	token      string
	checkPoint time.Time

	HTTPClient        *http.Client
	TokenLiveDuration time.Duration
}

// New return a new Mpesa Service
func New(key, secret string, endpoint string) *Service {
	if endpoint == "" {
		endpoint = SandboxEndpoint
	}
	b := []byte(key + ":" + secret)
	encoded := base64.StdEncoding.EncodeToString(b)
	serviceAuthHeader := "Basic " + encoded

	return &Service{
		appKey:            key,
		appSecret:         secret,
		endpoint:          endpoint,
		authHeader:        serviceAuthHeader,
		TokenLiveDuration: defaultTokenLive,
		HTTPClient:        http.DefaultClient,
	}
}

// Usually, service generate tokens on its own and you should not regenerate them manually.
func (s *Service) GenerateNewAccessToken() (string, error) {
	err := s.updateToken()
	if err != nil {
		return "", err
	}
	return s.token, nil
}

func (s *Service) checkToken() error {
	if len(s.token) == 0 || time.Since(s.checkPoint) > s.TokenLiveDuration {
		return ErrTokenIsExpired
	}
	return nil
}

// Generate Mpesa Daraja Access Token
func (s *Service) updateToken() error {
	url := s.endpoint + "oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add(authHeader, s.authHeader)

	client := s.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send auth request: %v", err)
	}
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	if resp.StatusCode != http.StatusOK {
		var apiErr errorResponse
		if err := dec.Decode(&apiErr); err == nil {
			return error(apiErr)
		}
		return errors.New(resp.Status)
	}

	var authResponse authResponse
	if err := dec.Decode(&authResponse); err != nil {
		return fmt.Errorf("could not decode auth response: %v", err)
	}

	s.token = authResponse.AccessToken
	if authResponse.ExpiresIn != "" {
		expSecs, err := strconv.Atoi(authResponse.ExpiresIn)
		if err == nil {
			s.TokenLiveDuration = time.Second * time.Duration(expSecs)
		}
	}
	s.checkPoint = time.Now()
	return nil
}

func (s *Service) roundTrip(reqBody interface{}, dest interface{}, url string) error {
	if s.checkToken() != nil {
		if err := s.updateToken(); err != nil {
			return fmt.Errorf("update auth token: %v", err)
		}
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("encode to json: %v", err)
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	r.Header.Add(authHeader, "Bearer "+s.token)
	r.Header.Add(contentTypeHeader, "application/json")

	resp, err := s.HTTPClient.Do(r)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()
	/*{
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
	}*/
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	if resp.StatusCode != http.StatusOK {
		var apiErr errorResponse
		if err := dec.Decode(&apiErr); err == nil {
			return error(apiErr)
		}
		return errors.New(resp.Status)
	}

	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("could not decode response: %v", err)
	}

	return nil
}

func (s *Service) C2BRegisterURL(c2BRegisterURL C2BRegisterURL) (*C2BRegisterURLResponse, error) {
	url := s.endpoint + "mpesa/c2b/v1/registerurl"
	var res C2BRegisterURLResponse
	err := s.roundTrip(c2BRegisterURL, &res, url)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) C2BSimulation(c2b C2B) (*C2BResponse, error) {
	url := s.endpoint + "mpesa/c2b/v1/simulate"
	var res C2BResponse
	err := s.roundTrip(c2b, &res, url)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) B2CRequest(b2c B2C) (*B2CResponse, error) {
	url := s.endpoint + "mpesa/b2c/v1/paymentrequest"
	var res B2CResponse
	err := s.roundTrip(b2c, &res, url)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) TransactionStatus(status TransactionStatus) (*TransactionStatusResponse, error) {
	url := s.endpoint + "mpesa/transactionstatus/v1/query"
	var res TransactionStatusResponse
	err := s.roundTrip(status, &res, url)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) MPESAOnlinePayment(payment Payment) (*PaymentResponse, error) {
	url := s.endpoint + "mpesa/stkpush/v1/processrequest"
	var res PaymentResponse
	err := s.roundTrip(payment, &res, url)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) Reversal(reversal Reversal) (*ReversalResponse, error) {
	url := s.endpoint + "mpesa/reversal/v1/request"
	var res ReversalResponse
	err := s.roundTrip(reversal, &res, url)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

/*
// B2BRequest sends a new request
func (s Service) B2BRequest(b2b B2B) (string, error) {
	body, err := json.Marshal(b2b)
	if err != nil {
		return "", nil
	}
	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/b2b/v1/paymentrequest"
	return s.newStringRequest(url, body, headers)
}

// BalanceInquiry sends a balance inquiry
func (s Service) BalanceInquiry(balanceInquiry BalanceInquiry) (string, error) {
	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	body, err := json.Marshal(balanceInquiry)
	if err != nil {
		return "", err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"
	headers["postman-token"] = "2aa448be-7d56-a796-065f-b378ede8b136"

	url := s.baseURL() + "mpesa/accountbalance/v1/query"
	return s.newStringRequest(url, body, headers)
}

func (s Service) newStringRequest(url string, body []byte, headers map[string]string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", nil
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}

	stringBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	log.Println("Response received")
	return string(stringBody), nil
}
*/
