package zota

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/deposit/shared"
)

var (
	logger         *zap.Logger
	depositGateway *DepositGateway
	requestPayload shared.Request
)

func setup(t *testing.T) {
	logger, _ = zap.NewDevelopment()

	cfg := &config.Config{
		ZotaBaseUrl:            "https://example.com",
		ZotaEndpointId:         "testEndpoint",
		ZotaAPISecretKey:       "testSecret",
		ZotaDepositCallBackUrl: "https://example.com/callback",
		ZotaDepositRedirectUrl: "https://example.com/redirect",
	}

	depositGateway = &DepositGateway{
		logger: logger,
		config: cfg,
	}

	requestPayload = shared.Request{
		ClientRequest: shared.ClientRequest{
			UserId:              "user123",
			OrderAmount:         "100.00",
			OrderCurrency:       "USD",
			CustomerEmail:       "test@example.com",
			CustomerFirstName:   "John",
			CustomerLastName:    "Doe",
			CustomerAddress:     "123 Main St",
			CustomerCountryCode: "US",
			CustomerCity:        "New York",
			CustomerState:       "NY",
			CustomerZipCode:     "10001",
			CustomerPhone:       "1234567890",
			CustomerIp:          "127.0.0.1",
			CheckoutUrl:         "https://example.com/checkout",
		},
	}
}

func TestDeposit_Success(t *testing.T) {
	setup(t)

	depositResponse := DepositResponse{
		Code: "200",
		Data: &DepositResponseData{
			DepositUrl:      "https://example.com/deposit",
			MerchantOrderID: uuid.New().String(),
			OrderID:         "123123",
		},
	}
	depositResponseJSON, _ := json.Marshal(depositResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/deposit/request/testEndpoint/", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write(depositResponseJSON)
	}))
	defer server.Close()

	depositGateway.config.ZotaBaseUrl = server.URL

	response, err := depositGateway.Deposit(requestPayload)
	require.NoError(t, err)
	assert.Equal(t, depositResponse.Data.MerchantOrderID, response.OrderID)
	assert.Equal(t, depositResponse.Data.OrderID, response.PaymentGatewayOrderID)
}

func TestDeposit_BuildDepositReqError(t *testing.T) {
	setup(t)

	badRequest := shared.Request{}
	response, err := depositGateway.Deposit(badRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestDeposit_HTTPRequestError(t *testing.T) {
	setup(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "network error", http.StatusInternalServerError)
	}))
	defer server.Close()

	depositGateway.config.ZotaBaseUrl = server.URL

	response, err := depositGateway.Deposit(requestPayload)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestDeposit_NonOKHTTPResponse(t *testing.T) {
	setup(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	depositGateway.config.ZotaBaseUrl = server.URL

	response, err := depositGateway.Deposit(requestPayload)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestDeposit_UnmarshalResponseError(t *testing.T) {
	setup(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid JSON"))
	}))
	defer server.Close()

	depositGateway.config.ZotaBaseUrl = server.URL

	response, err := depositGateway.Deposit(requestPayload)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestBuildDepositReq(t *testing.T) {
	setup(t)

	depositReq, err := depositGateway.buildDepositReq(requestPayload)
	require.NoError(t, err)
	assert.Equal(t, "Deposit", depositReq.MerchantOrderDesc)
	assert.Equal(t, "100.00", depositReq.OrderAmount)
	assert.Equal(t, "USD", depositReq.OrderCurrency)
	assert.Equal(t, "test@example.com", depositReq.CustomerEmail)
	assert.NotEmpty(t, depositReq.Signature)
}

func TestMarshalCustomParam(t *testing.T) {
	setup(t)

	userId := "user123"
	customParamJSON, err := depositGateway.marshalCustomParam(userId)
	require.NoError(t, err)
	assert.Contains(t, customParamJSON, `"UserId":"user123"`)
}

func TestBuildSignature(t *testing.T) {
	setup(t)

	orderAmount := "100.00"
	customerEmail := "test@example.com"
	endpointID := "testEndpoint"
	merchantOrderId := uuid.New()
	merchantSecretKey := "testSecret"

	signature := depositGateway.buildSignature(orderAmount, customerEmail, endpointID, merchantOrderId, merchantSecretKey)
	assert.NotEmpty(t, signature)

	expectedSignatureString := fmt.Sprintf("%s%s%s%s%s", endpointID, merchantOrderId, orderAmount, customerEmail, merchantSecretKey)
	hash := sha256.New()
	hash.Write([]byte(expectedSignatureString))
	expectedSignature := hex.EncodeToString(hash.Sum(nil))
	assert.Equal(t, expectedSignature, signature)
}

func TestHandleDepositResponse(t *testing.T) {
	setup(t)

	depositResponse := DepositResponse{
		Code: "200",
		Data: &DepositResponseData{
			DepositUrl:      "https://example.com/deposit",
			MerchantOrderID: uuid.New().String(),
			OrderID:         uuid.New().String(),
		},
	}
	depositResponseJSON, _ := json.Marshal(depositResponse)

	response, err := depositGateway.handleDepositResponse(depositResponseJSON, requestPayload)
	require.NoError(t, err)
	assert.Equal(t, depositResponse.Data.MerchantOrderID, response.OrderID)
	assert.Equal(t, depositResponse.Data.OrderID, response.PaymentGatewayOrderID)
}
