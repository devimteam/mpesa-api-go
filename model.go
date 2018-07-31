package mpesa

import "fmt"

type authResponse struct {
	RequestId    *string `json:"requestId"`
	ErrorCode    *string `json:"errorCode"`
	ErrorMessage *string `json:"errorMessage"`

	// Access token to access other APIs
	AccessToken string `json:"access_token"`
	// Token expiry time in seconds.
	ExpiresIn string `json:"expires_in"`
}

func (r authResponse) Error() string {
	return fmt.Sprintf("code %s - %s", sp(r.ErrorCode), sp(r.ErrorMessage))
}

func sp(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// C2B is a model
type C2B struct {
	ShortCode     string
	CommandID     string
	Amount        string
	Msisdn        string
	BillRefNumber string
}

// B2C is a model
type B2C struct {
	InitiatorName      string
	SecurityCredential string
	CommandID          string
	Amount             string
	PartyA             string
	PartyB             string
	Remarks            string
	QueueTimeOutURL    string
	ResultURL          string
	Occassion          string
}

// B2B is a model
type B2B struct {
	Initiator              string
	SecurityCredential     string
	CommandID              string
	SenderIdentifierType   string
	RecieverIdentifierType string
	Amount                 string
	PartyA                 string
	PartyB                 string
	Remarks                string
	AccountReference       string
	QueueTimeOutURL        string
	ResultURL              string
}

// STKPush is a model
type MPESAExpress struct {
	BusinessShortCode string
	Password          string
	Timestamp         string
	TransactionType   string
	Amount            string
	PartyA            string
	PartyB            string
	PhoneNumber       string
	CallBackURL       string
	AccountReference  string
	TransactionDesc   string
}

// Reversal is a model
type Reversal struct {
	Initiator              string
	SecurityCredential     string
	CommandID              string
	TransactionID          string
	Amount                 string
	ReceiverParty          string
	ReceiverIdentifierType string
	QueueTimeOutURL        string
	ResultURL              string
	Remarks                string
	Occassion              string
}

// BalanceInquiry is a model
type BalanceInquiry struct {
	Initiator          string
	SecurityCredential string
	CommandID          string
	PartyA             string
	IdentifierType     string
	Remarks            string
	QueueTimeOutURL    string
	ResultURL          string
}

// RegisterURL is a model
type C2BRegisterURL struct {
	// The short code of the organization.
	ShortCode string
	// Default response type for timeout. Incase a tranaction times out, Mpesa will by default Complete or Cancel the transaction.
	ResponseType string
	// Confirmation URL for the client.
	ConfirmationURL string
	// Validation URL for the client.
	ValidationURL string
}

type C2BRegisterURLResponse struct {
	OriginatorConverstionID string `json:"originatorConverstionId"`
	ConversationID          string `json:"conversationId"`
	ResponseDescription     string `json:"responseDescription"`
}
