package zota

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/status/shared"
)

type statusGatewayTestSuite struct {
	mockCtrl      *gomock.Controller
	logger        *zap.Logger
	config        *config.Config
	statusGateway *StatusGateway
	request       shared.Request
}

func (s *statusGatewayTestSuite) setup(t *testing.T) {
	s.mockCtrl = gomock.NewController(t)
	s.logger, _ = zap.NewDevelopment()

	s.config = &config.Config{
		ZotaBaseUrl:            "https://example.com",
		ZotaMerchantId:         "merchant123",
		ZotaAPISecretKey:       "secret123",
		ZotaEndpointId:         "endpoint123",
		ZotaDepositCallBackUrl: "https://callback.example.com",
		ZotaDepositRedirectUrl: "https://redirect.example.com",
	}

	s.statusGateway = NewStatusGateway(s.logger, s.config)

	s.request = shared.Request{
		ClientRequest: shared.ClientRequest{
			OrderId:         "order123",
			MerchantOrderId: "merchantOrder123",
		},
	}
}

func (s *statusGatewayTestSuite) teardown() {
	s.mockCtrl.Finish()
}

func TestCheckStatus_Success(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	statusResponse := StatusResponse{
		Code: "200",
		Data: &Data{
			Type:          "type1",
			Status:        "status1",
			Amount:        "100.00",
			Currency:      "USD",
			CustomerEmail: "test@example.com",
		},
	}
	statusResponseJSON, _ := json.Marshal(statusResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/query/order-status/", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write(statusResponseJSON)
	}))
	defer server.Close()

	s.statusGateway.config.ZotaBaseUrl = server.URL

	response, err := s.statusGateway.CheckStatus(s.request)
	require.NoError(t, err)
	assert.Equal(t, "type1", response.Type)
	assert.Equal(t, "status1", response.Status)
	assert.Equal(t, "100.00", response.Amount)
	assert.Equal(t, "USD", response.Currency)
	assert.Equal(t, "test@example.com", response.CustomerEmail)
}

func TestCheckStatus_BuildStatusReqError(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	// Simulate a failure in building the status request by using invalid configuration
	s.statusGateway.config.ZotaMerchantId = ""

	response, err := s.statusGateway.CheckStatus(s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestCheckStatus_HTTPRequestError(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	// Simulate a server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "network error", http.StatusInternalServerError)
	}))
	defer server.Close()

	s.statusGateway.config.ZotaBaseUrl = server.URL

	response, err := s.statusGateway.CheckStatus(s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestCheckStatus_NonOKHTTPResponse(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	// Simulate a server that returns a non-OK response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	s.statusGateway.config.ZotaBaseUrl = server.URL

	response, err := s.statusGateway.CheckStatus(s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestCheckStatus_UnmarshalResponseError(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	// Simulate a server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid JSON"))
	}))
	defer server.Close()

	s.statusGateway.config.ZotaBaseUrl = server.URL

	response, err := s.statusGateway.CheckStatus(s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestBuildStatusReq(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	statusReq, err := s.statusGateway.buildStatusReq(s.request)
	require.NoError(t, err)
	assert.Equal(t, s.config.ZotaMerchantId, statusReq.MerchantId)
	assert.Equal(t, s.request.OrderId, statusReq.OrderId)
	assert.Equal(t, s.request.MerchantOrderId, statusReq.MerchantOrderId)
	assert.NotEmpty(t, statusReq.Timestamp)
	assert.NotEmpty(t, statusReq.Signature)
}

func TestBuildSignature(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := s.statusGateway.buildSignature(s.request, timestamp)
	assert.NotEmpty(t, signature)

	expectedSignatureString := fmt.Sprintf("%s%s%s%s%s", s.config.ZotaMerchantId, s.request.MerchantOrderId, s.request.OrderId, timestamp, s.config.ZotaAPISecretKey)
	hash := sha256.New()
	hash.Write([]byte(expectedSignatureString))
	expectedSignature := hex.EncodeToString(hash.Sum(nil))
	assert.Equal(t, expectedSignature, signature)
}

func TestHandleStatusResponse(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	statusResponse := StatusResponse{
		Code: "200",
		Data: &Data{
			Type:          "type1",
			Status:        "status1",
			Amount:        "100.00",
			Currency:      "USD",
			CustomerEmail: "test@example.com",
		},
	}
	statusResponseJSON, _ := json.Marshal(statusResponse)

	response, err := s.statusGateway.handleStatusResponse(statusResponseJSON, s.request)
	require.NoError(t, err)
	assert.Equal(t, "type1", response.Type)
	assert.Equal(t, "status1", response.Status)
	assert.Equal(t, "100.00", response.Amount)
	assert.Equal(t, "USD", response.Currency)
	assert.Equal(t, "test@example.com", response.CustomerEmail)
}

func TestHandleStatusResponse_NoDataError(t *testing.T) {
	s := &statusGatewayTestSuite{}
	s.setup(t)
	defer s.teardown()

	statusResponse := StatusResponse{
		Code: "200",
		Data: nil,
	}
	statusResponseJSON, _ := json.Marshal(statusResponse)

	response, err := s.statusGateway.handleStatusResponse(statusResponseJSON, s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "no data in statusResponse", err.Error())
}
