package shared

/*
// ClientRequest - represents the expected request body for the deposit endpoint in the merchant server
// for the sake of simplicity we will accept all the required fields as req body
// in real-world scenario we will not receive most of the data and
// we can not trust the data received from the client
// Also , let's pretend we have userId for the user
type ClientRequest struct {
	UserId        int    `json:"userId"`
	OrderAmount   string `json:"orderAmount"`
	OrderCurrency string `json:"orderCurrency"`
	CustomerEmail string `json:"customerEmail"`
}
*/

type ClientRequest struct {
	UserId              string `json:"userId" validate:"required"`
	OrderAmount         string `json:"orderAmount" validate:"required"`
	OrderCurrency       string `json:"orderCurrency" validate:"required"`
	CustomerEmail       string `json:"customerEmail" validate:"required,email"`
	CustomerFirstName   string `json:"customerFirstName" validate:"required"`
	CustomerLastName    string `json:"customerLastName" validate:"required"`
	CustomerAddress     string `json:"customerAddress" validate:"required"`
	CustomerCountryCode string `json:"customerCountryCode" validate:"required"`
	CustomerCity        string `json:"customerCity" validate:"required"`
	CustomerZipCode     string `json:"customerZipCode" validate:"required"`
	CustomerPhone       string `json:"customerPhone" validate:"required"`
	CustomerIp          string `json:"customerIp" validate:"required"`
	CheckoutUrl         string `json:"checkoutUrl" validate:"url"`
	Language            string `json:"language"`
	CustomerState       string `json:"customerState"`
	CustomerBankCode    string `json:"customerBankCode"`
}

// Request - service level model
type Request struct {
	ClientRequest
}

type Response struct {
	ClientRequest         ClientRequest `json:"request"`
	OrderID               string        `json:"orderId"`
	PaymentGatewayOrderID string        `json:"paymentGatewayOrderId"`
}

type DepositPaymentGateway interface {
	Deposit(req Request) (*Response, error)
}
