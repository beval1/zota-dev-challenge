package zota

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/status/shared"
)

const StatusCheckApiPath = "api/v1/query/order-status"

type StatusRequest struct {
	MerchantId      string `json:"merchantID" schema:"merchantID"`
	OrderId         string `json:"orderID" schema:"orderID"`
	MerchantOrderId string `json:"merchantOrderID" schema:"merchantOrderID"`
	Timestamp       string `json:"timestamp" schema:"timestamp"`
	Signature       string `json:"signature" schema:"signature"`
}

type StatusResponse struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

type Data struct {
	Type                   string        `json:"type"`
	Status                 string        `json:"status"`
	ErrorMessage           string        `json:"errorMessage,omitempty"`
	EndpointID             string        `json:"endpointID"`
	ProcessorTransactionID string        `json:"processorTransactionID"`
	OrderID                string        `json:"orderID"`
	MerchantOrderID        string        `json:"merchantOrderID"`
	Amount                 string        `json:"amount"`
	Currency               string        `json:"currency"`
	CustomerEmail          string        `json:"customerEmail"`
	CustomParam            string        `json:"customParam"`
	ExtraData              ExtraData     `json:"extraData"`
	Request                StatusRequest `json:"request"`
}

type ExtraData struct {
	AmountChanged     bool   `json:"amountChanged"`
	AmountRounded     bool   `json:"amountRounded"`
	AmountManipulated bool   `json:"amountManipulated"`
	Dcc               bool   `json:"dcc"`
	OriginalAmount    string `json:"originalAmount"`
	PaymentMethod     string `json:"paymentMethod"`
	SelectedBankCode  string `json:"selectedBankCode"`
	SelectedBankName  string `json:"selectedBankName"`
}

type StatusGateway struct {
	logger *zap.Logger
	config *config.Config
}

func NewStatusGateway(logger *zap.Logger, config *config.Config) *StatusGateway {
	return &StatusGateway{logger: logger, config: config}
}

func (s *StatusGateway) CheckStatus(req shared.Request) (*shared.Response, error) {
	s.logger.Info("Checking status with Zota", zap.Any("request", req))

	statusReq, err := s.buildStatusReq(req)
	if err != nil {
		s.logger.Error("Failed to build status request", zap.Error(err))
		return nil, err
	}

	statusCheckApiUrl, err := s.buildStatusCheckApiUrl(statusReq)
	if err != nil {
		s.logger.Error("Failed to build status check API URL", zap.Error(err))
		return nil, err
	}

	respBody, statusCode, err := s.sendStatusRequest(statusCheckApiUrl)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK statusResponse from Zota server: %s", http.StatusText(statusCode))
	}

	response, err := s.handleStatusResponse(respBody, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *StatusGateway) buildStatusReq(req shared.Request) (StatusRequest, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := s.buildSignature(req, timestamp)

	return StatusRequest{
		MerchantId:      s.config.ZotaMerchantId,
		OrderId:         req.OrderId,
		MerchantOrderId: req.MerchantOrderId,
		Timestamp:       timestamp,
		Signature:       signature,
	}, nil
}

func (s *StatusGateway) buildSignature(req shared.Request, timestamp string) string {
	signatureString := fmt.Sprintf("%s%s%s%s%s", s.config.ZotaMerchantId, req.MerchantOrderId, req.OrderId, timestamp, s.config.ZotaAPISecretKey)
	hash := sha256.New()
	hash.Write([]byte(signatureString))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *StatusGateway) buildStatusCheckApiUrl(statusReq StatusRequest) (string, error) {
	encoder := schema.NewEncoder()

	values := url.Values{}
	if err := encoder.Encode(statusReq, values); err != nil {
		return "", fmt.Errorf("failed to marshal request to query parameters: %w", err)
	}

	return fmt.Sprintf("%s/%s/?%s", s.config.ZotaBaseUrl, StatusCheckApiPath, values.Encode()), nil
}

func (s *StatusGateway) sendStatusRequest(statusCheckApiUrl string) ([]byte, int, error) {
	s.logger.Info("Sending status request to Zota server", zap.String("url", statusCheckApiUrl))

	httpReq, err := http.NewRequest("GET", statusCheckApiUrl, nil)
	if err != nil {
		s.logger.Error("Failed to create HTTP request", zap.Error(err))
		return nil, 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		s.logger.Error("Failed to send HTTP request", zap.Error(err))
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read statusResponse body", zap.Error(err))
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

func (s *StatusGateway) handleStatusResponse(respBody []byte, req shared.Request) (*shared.Response, error) {
	var statusResponse StatusResponse
	if err := json.Unmarshal(respBody, &statusResponse); err != nil {
		s.logger.Error("Failed to unmarshal statusResponse body", zap.Error(err))
		return nil, err
	}

	if statusResponse.Data == nil {
		return nil, fmt.Errorf("no data in statusResponse")
	}

	response := shared.Response{
		ClientRequest: req.ClientRequest,
		Type:          statusResponse.Data.Type,
		Status:        statusResponse.Data.Status,
		Amount:        statusResponse.Data.Amount,
		Currency:      statusResponse.Data.Currency,
		CustomerEmail: statusResponse.Data.CustomerEmail,
	}

	s.logger.Info("Successfully checked status request",
		zap.Any("request", req),
		zap.ByteString("statusResponse", respBody))
	return &response, nil
}
