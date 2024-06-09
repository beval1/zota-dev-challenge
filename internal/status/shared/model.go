package shared

type ClientRequest struct {
	OrderId         string `json:"orderId"`
	MerchantOrderId string `json:"merchantOrderId"`
}

type Request struct {
	ClientRequest
}

type Response struct {
	ClientRequest ClientRequest `json:"request"`
	Type          string        `json:"type"`
	Status        string        `json:"status"`
	Amount        string        `json:"amount"`
	Currency      string        `json:"currency"`
	CustomerEmail string        `json:"customerEmail"`
}

type StatusPaymentGateway interface {
	CheckStatus(req Request) (*Response, error)
}
