package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"zota-dev-challenge/internal/deposit/shared"
)

var (
	mockCtrl       *gomock.Controller
	mockService    *MockServiceInterface
	logger         *zap.Logger
	validate       *validator.Validate
	requestPayload shared.ClientRequest
)

func setup(t *testing.T) {
	mockCtrl = gomock.NewController(t)
	mockService = NewMockServiceInterface(mockCtrl)
	logger, _ = zap.NewDevelopment()
	validate = validator.New()

	// Set a common request payload
	requestPayload = shared.ClientRequest{
		UserId:              "333",
		OrderAmount:         "1.0",
		OrderCurrency:       "USD",
		CustomerEmail:       "john@gmail.com",
		CustomerFirstName:   "John",
		CustomerLastName:    "Lock",
		CustomerAddress:     "Sofia",
		CustomerCountryCode: "US",
		CustomerCity:        "Sofia",
		CustomerState:       "FL",
		CustomerZipCode:     "32042",
		CustomerPhone:       "+1 420-100-1000",
		CustomerIp:          "192.168.1.1",
		CheckoutUrl:         "http://example.com",
	}
}

func teardown() {
	mockCtrl.Finish()
}

func TestHandler(t *testing.T) {
	setup(t)
	defer teardown()

	mockService.EXPECT().
		ProcessDeposit(gomock.Any()).
		Return(&shared.Response{requestPayload, "1", "11"}, nil)

	handler := Handler(mockService, logger, validate)

	payloadBytes, _ := json.Marshal(requestPayload)
	req, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response shared.Response
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, requestPayload, response.ClientRequest)
}

func TestHandler_InvalidJSON(t *testing.T) {
	setup(t)
	defer teardown()

	handler := Handler(mockService, logger, validate)

	req, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_MissingRequiredFields(t *testing.T) {
	setup(t)
	defer teardown()

	requestPayload.CustomerEmail = ""

	handler := Handler(mockService, logger, validate)

	payloadBytes, _ := json.Marshal(requestPayload)
	req, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_InvalidFieldValues(t *testing.T) {
	setup(t)
	defer teardown()

	requestPayload.CustomerEmail = "invalid-email"

	handler := Handler(mockService, logger, validate)

	payloadBytes, _ := json.Marshal(requestPayload)
	req, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_ServiceError(t *testing.T) {
	setup(t)
	defer teardown()

	mockService.EXPECT().
		ProcessDeposit(gomock.Any()).
		Return(nil, fmt.Errorf("service error"))

	handler := Handler(mockService, logger, validate)

	payloadBytes, _ := json.Marshal(requestPayload)
	req, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
