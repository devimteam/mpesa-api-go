// Every comment was copied from original documentation from https://developer.safaricom.co.ke/apis-explorer
package mpesa

import (
	"encoding/json"
	"fmt"
)

type authResponse struct {
	// Access token to access other APIs
	AccessToken string `json:"access_token"`
	// Token expiry time in seconds.
	ExpiresIn string `json:"expires_in"`
}

type errorResponse struct {
	RequestId    *string `json:"requestId"`
	ErrorCode    *string `json:"errorCode"`
	ErrorMessage *string `json:"errorMessage"`
}

func (r errorResponse) Error() string {
	return fmt.Sprintf("%s - %s", sp(r.ErrorCode), sp(r.ErrorMessage))
}

func sp(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

type GenericResponse struct {
	// A unique numeric code generated by the M-Pesa system of the request.
	OriginatorConversationID string `json:",omitempty"`
	// A unique numeric code generated by the M-Pesa system of the response to a request.
	ConversationID      string `json:",omitempty"`
	ResponseDescription string `json:",omitempty"`
	// A response message from the M-Pesa system accompanying the response to a request.
	ResponseCode string `json:",omitempty"`
}

type C2B struct {
	// Short Code receiving the amount being transacted
	ShortCode string
	// Unique command for each transaction type. For C2B dafult
	// CustomerPayBillOnline
	// CustomerBuyGoodsOnline
	CommandID string
	// The amount being transacted
	Amount string
	// Phone number (msisdn) initiating the transaction
	Msisdn string
	// Bill Reference Number (Optional)
	BillRefNumber string `json:",omitempty"`
}

type C2BResponse GenericResponse

type B2C struct {
	// The name of the initiator initiating the request
	// This is the credential/username used to authenticate the transaction request
	InitiatorName string
	// Encrypted Credential of user getting transaction amount
	// Encrypted password for the initiator to authenticate the transaction request
	SecurityCredential string
	// Unique command for each transaction type
	// SalaryPayment
	// BusinessPayment
	// PromotionPayment
	CommandID string
	// The amount been transacted
	Amount string
	// Organization /MSISDN sending the transaction
	// Shortcode (6 digits)
	// MSISDN (12 digits)
	PartyA string
	// MSISDN receiving the transaction (12 digits)
	PartyB string
	// Comments that are sent along with the transaction.
	// Up to 100
	Remarks string
	// The path that stores information of time out transaction
	// https://ip or domain:port/path
	QueueTimeOutURL string
	// The path that stores information of transactions
	// https://ip or domain:port/path
	ResultURL string
	// Optional Parameter
	// Up to 100
	Occasion string `json:",omitempty"`
}

type B2CResponse GenericResponse

type B2CCallback struct {
	Result struct {
		ResultType               int
		ResultCode               int
		ResultDesc               string
		OriginatorConversationID string
		ConversationID           string
		TransactionID            string
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string
				Value json.RawMessage
			}
		}
		ReferenceData json.RawMessage
	}
}

type TransactionStatus struct {
	// Takes only 'TransactionStatusQuery' command id
	CommandID string
	// Organization/MSISDN receiving the transaction
	// -Shortcode (6 digits)
	// -MSISDN (12 Digits)
	PartyA string
	// Type of organization receiving the transaction
	IdentifierType string
	// Comments that are sent along with the transaction
	// Up to 100
	Remarks string
	// The name of Initiator to initiating  the request
	// This is the credential/username used to authenticate the transaction request
	Initiator string
	// Encrypted Credential of user getting transaction amount
	// Encrypted password for the initiator to authenticate the transaction request
	SecurityCredential string
	// The path that stores information of time out transaction
	// https://ip or domain:port/path
	QueueTimeOutURL string
	// The path that stores information of transaction
	// https://ip or domain:port/path
	ResultURL string
	// Unique identifier to identify a transaction on M-Pesa
	TransactionID string
	Occasion      string `json:",omitempty"`
}

type TransactionStatusResponse GenericResponse

// https://developer.safaricom.co.ke/lipa-na-m-pesa-online/apis/post/stkpush/v1/processrequest
type Payment struct {
	// This is organizations shortcode (Paybill or Buygoods - A 5 to 6 digit account number)
	// used to identify an organization and receive the transaction.
	BusinessShortCode string
	// This is the password used for encrypting the request sent: A base64 encoded string.
	// (The base64 string is a combination of Shortcode+Passkey+Timestamp)
	Password string
	// This is the Timestamp of the transaction,
	// normally in the format of YEAR+MONTH+DATE+HOUR+MINUTE+SECOND (YYYYMMDDHHMMSS)
	// Each part should be at least two digits apart from the year which takes four digits.
	Timestamp string
	// This is the transaction type that is used to identify the transaction when sending the request to M-Pesa.
	// The transaction type for M-Pesa Express is "CustomerPayBillOnline"
	TransactionType string
	// This is the Amount transacted normally a numeric value. Money that customer pays to the Shorcode.
	// Only whole numbers are supported.
	Amount string
	// The phone number sending money.
	// The parameter expected is a Valid Safaricom Mobile Number that is M-Pesa registered in the format 2547XXXXXXXX
	PartyA string
	// The organization receiving the funds.
	// The parameter expected is a 5 to 6 digit as defined on the Shortcode description above.
	// This can be the same as BusinessShortCode value above.
	PartyB string
	// The Mobile Number to receive the STK Pin Prompt. This number can be the same as PartyA value above.
	PhoneNumber string
	// A CallBack URL is a valid secure URL that is used to receive notifications from M-Pesa API.
	// It is the endpoint to which the results will be sent by M-Pesa API.
	CallBackURL string
	// Account Reference:
	// This is an Alpha-Numeric parameter that is defined by your system as an Identifier of the transaction
	// for CustomerPayBillOnline transaction type.
	// Along with the business name, this value is also displayed to the customer in the STK Pin Prompt message.
	// Maximum of 12 characters.
	AccountReference string
	// This is any additional information/comment that can be sent along with the request from your system.
	// Maximum of 13 Characters.
	TransactionDesc string
}

type PaymentResponse struct {
	// This is a global unique Identifier for any submitted payment request.
	MerchantRequestID string
	// This is a global unique identifier of the processed checkout transaction request.
	CheckoutRequestID string
	// Response description is an acknowledgement message from the API that gives the status of the request submission,
	// usually maps to a specific ResponseCode value.
	// It can be a Success submission message or an error description.
	ResponseDescription string
	// This is a Numeric status code that indicates the status of the transaction submission.
	// 0 means successful submission and any other code means an error occurred.
	ResponseCode string
	// This is a message that your system can display to the Customer as an acknowledgement of the payment request submission.
	CustomerMessage string
}

type PaymentCallback struct {
	Body struct {
		STKCallback struct {
			// This is a global unique Identifier for any submitted payment request.
			// This is the same value returned in the acknowledgement message of the initial request.
			MerchantRequestID string
			// This is a global unique identifier of the processed checkout transaction request.
			// This is the same value returned in the acknowledgement message of the initial request.
			CheckoutRequestID string
			// This is a numeric status code that indicates the status of the transaction processing.
			// 0 means successful processing and any other code means an error occurred or the transaction failed.
			ResultCode string
			// Result description is a message from the API that gives the status of the request processing,
			// usually maps to a specific ResultCode value.
			// It can be a Success processing message or an error description message.
			ResultDesc       string
			CallbackMetadata struct {
				Item []struct {
					Name  string
					Value json.RawMessage `json:",omitempty"`
				}
			} `json:",omitempty"`
		} `json:"stkCallback"`
	}
}

type C2BRegisterURL struct {
	// The short code of the organization.
	ShortCode string
	// Default response type for timeout. In case a transaction times out,
	// Mpesa will by default Complete or Cancel the transaction.
	ResponseType string
	// Confirmation URL for the client.
	ConfirmationURL string
	// Validation URL for the client.
	ValidationURL string
}

type C2BRegisterURLResponse GenericResponse

type Reversal struct {
	// Takes only 'TransactionReversal' Command id
	CommandID string
	// Organization receiving the transaction (shortcode)
	ReceiverParty string
	// Type of organization receiving the transaction
	// Organization Identifier on M-Pesa
	ReceiverIdentifierType string
	// Comments that are sent along with the transaction.
	// Up to 100 characters.
	Remarks string
	// The name of Initiator to initiating  the request
	// This is the credential/username used to authenticate the transaction request
	Initiator string
	// Encrypted Credential of user getting transaction amount
	// Encrypted password for the initiator to authenticate the transaction request
	SecurityCredential string
	// The path that stores information of time out transaction
	// https://ip or domain:port/path
	QueueTimeOutURL string
	// Organization Receiving the funds // WTF??
	TransactionID string
	// Optional Parameter
	// Up to 100 characters
	Occasion string
}

type ReversalResponse GenericResponse

/*
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
*/
