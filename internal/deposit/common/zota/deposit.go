package zota

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/deposit/shared"
)

const PaymentGatewayDepositApiPath = "api/v1/deposit/request"

type DepositRequest struct {
	MerchantOrderID     string `json:"merchantOrderID"`
	MerchantOrderDesc   string `json:"merchantOrderDesc"`
	OrderAmount         string `json:"orderAmount"`
	OrderCurrency       string `json:"orderCurrency"`
	CustomerEmail       string `json:"customerEmail"`
	CustomerFirstName   string `json:"customerFirstName"`
	CustomerLastName    string `json:"customerLastName"`
	CustomerAddress     string `json:"customerAddress"`
	CustomerCountryCode string `json:"customerCountryCode"`
	CustomerCity        string `json:"customerCity"`
	CustomerState       string `json:"customerState"`
	CustomerZipCode     string `json:"customerZipCode"`
	CustomerPhone       string `json:"customerPhone"`
	CustomerBankCode    string `json:"customerBankCode"`
	CustomerIP          string `json:"customerIP"`
	RedirectUrl         string `json:"redirectUrl"`
	CallbackUrl         string `json:"callbackUrl"`
	CustomParam         string `json:"customParam"`
	CheckoutUrl         string `json:"checkoutUrl"`
	Signature           string `json:"signature"`
}

type CustomParam struct {
	UserId string `json:"UserId"`
}

type DepositResponse struct {
	Code    string               `json:"code"`
	Message string               `json:"message,omitempty"`
	Data    *DepositResponseData `json:"data,omitempty"`
}

type DepositResponseData struct {
	DepositUrl      string `json:"depositUrl"`
	MerchantOrderID string `json:"merchantOrderID"`
	OrderID         string `json:"orderID"`
}

type DepositGateway struct {
	logger *zap.Logger
	config *config.Config
}

func NewDepositGateway(logger *zap.Logger, config *config.Config) *DepositGateway {
	return &DepositGateway{logger: logger, config: config}
}

func (d *DepositGateway) Deposit(req shared.Request) (*shared.Response, error) {
	d.logger.Info("Processing deposit with Zota", zap.Any("request", req))

	depositReq, err := d.buildDepositReq(req)
	if err != nil {
		d.logger.Error("Failed to build deposit request", zap.Error(err))
		return nil, err
	}

	depositReqJSON, err := json.Marshal(depositReq)
	if err != nil {
		d.logger.Error("Failed to marshal deposit request", zap.Error(err))
		return nil, err
	}

	respBody, statusCode, err := d.sendDepositRequest(depositReqJSON)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK depositResponse from Zota server: %s", http.StatusText(statusCode))
	}

	response, err := d.handleDepositResponse(respBody, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DepositGateway) buildDepositReq(req shared.Request) (DepositRequest, error) {
	endpointID := d.config.ZotaEndpointId
	merchantSecretKey := d.config.ZotaAPISecretKey
	zotaDepositCallBackUrl := d.config.ZotaDepositCallBackUrl
	zotaDepositRedirectUrl := d.config.ZotaDepositRedirectUrl
	merchantOrderId := uuid.New()

	signature := d.buildSignature(req.OrderAmount, req.CustomerEmail, endpointID, merchantOrderId, merchantSecretKey)

	customParamJSON, err := d.marshalCustomParam(req.ClientRequest.UserId)
	if err != nil {
		return DepositRequest{}, err
	}

	return DepositRequest{
		MerchantOrderID:     merchantOrderId.String(),
		MerchantOrderDesc:   "Deposit",
		OrderAmount:         req.OrderAmount,
		OrderCurrency:       req.OrderCurrency,
		CustomerEmail:       req.CustomerEmail,
		CustomerFirstName:   req.CustomerFirstName,
		CustomerLastName:    req.CustomerLastName,
		CustomerAddress:     req.CustomerAddress,
		CustomerCountryCode: req.CustomerCountryCode,
		CustomerCity:        req.CustomerCity,
		CustomerState:       req.CustomerState,
		CustomerZipCode:     req.CustomerZipCode,
		CustomerPhone:       req.CustomerPhone,
		CustomerBankCode:    req.CustomerBankCode,
		CustomerIP:          req.CustomerIp,
		RedirectUrl:         zotaDepositRedirectUrl,
		CallbackUrl:         zotaDepositCallBackUrl,
		CustomParam:         customParamJSON,
		CheckoutUrl:         req.ClientRequest.CheckoutUrl,
		Signature:           signature,
	}, nil
}

func (d *DepositGateway) marshalCustomParam(userId string) (string, error) {
	customParam := CustomParam{UserId: userId}
	customParamJSON, err := json.Marshal(customParam)
	if err != nil {
		d.logger.Error("Failed to marshal custom param", zap.Error(err))
		return "", err
	}
	return string(customParamJSON), nil
}

func (d *DepositGateway) buildSignature(orderAmount, customerEmail, endpointID string, merchantOrderId uuid.UUID, merchantSecretKey string) string {
	signatureString := fmt.Sprintf("%s%s%s%s%s", endpointID, merchantOrderId, orderAmount, customerEmail, merchantSecretKey)
	hash := sha256.New()
	hash.Write([]byte(signatureString))
	return hex.EncodeToString(hash.Sum(nil))
}

func (d *DepositGateway) sendDepositRequest(depositReqJSON []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s/%s/", d.config.ZotaBaseUrl, PaymentGatewayDepositApiPath, d.config.ZotaEndpointId)
	d.logger.Debug("Sending deposit request to Zota server", zap.String("url", url))

	reqBody := bytes.NewBuffer(depositReqJSON)
	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		d.logger.Error("Failed to create HTTP request", zap.Error(err))
		return nil, 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		d.logger.Error("Failed to send HTTP request", zap.Error(err))
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		d.logger.Error("Failed to read depositResponse body", zap.Error(err))
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

func (d *DepositGateway) handleDepositResponse(respBody []byte, req shared.Request) (*shared.Response, error) {
	var depositResponse DepositResponse
	if err := json.Unmarshal(respBody, &depositResponse); err != nil {
		d.logger.Error("Failed to unmarshal depositResponse body", zap.Error(err))
		return nil, err
	}

	response := shared.Response{
		ClientRequest:         req.ClientRequest,
		OrderID:               depositResponse.Data.MerchantOrderID,
		PaymentGatewayOrderID: depositResponse.Data.OrderID,
	}

	d.logger.Info("Successfully processed deposit request",
		zap.Any("request", req),
		zap.ByteString("depositResponse", respBody))
	return &response, nil
}
