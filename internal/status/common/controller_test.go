package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"zota-dev-challenge/internal/status/shared"
)

type controllerTestSuite struct {
	mockCtrl    *gomock.Controller
	mockService *MockServiceInterface
	logger      *zap.Logger
	request     shared.ClientRequest
}

func (s *controllerTestSuite) setup(t *testing.T) {
	s.mockCtrl = gomock.NewController(t)
	s.mockService = NewMockServiceInterface(s.mockCtrl)
	s.logger, _ = zap.NewDevelopment()

	s.request = shared.ClientRequest{
		OrderId:         "order123",
		MerchantOrderId: "merchantOrder123",
	}
}

func (s *controllerTestSuite) teardown() {
	s.mockCtrl.Finish()
}

func TestHandler_Success(t *testing.T) {
	s := &controllerTestSuite{}
	s.setup(t)
	defer s.teardown()

	expectedResponse := &shared.Response{
		Status: "Success",
	}

	s.mockService.EXPECT().CheckStatus(&s.request).Return(expectedResponse, nil)

	request, err := http.NewRequest("GET", "/status?orderId=order123&merchantOrderId=merchantOrder123", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := Handler(s.mockService, s.logger)
	handler.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response shared.Response
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, &response)
}

func TestHandler_DecodeError(t *testing.T) {
	s := &controllerTestSuite{}
	s.setup(t)
	defer s.teardown()

	request, err := http.NewRequest("GET", "/status?invalidparam", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := Handler(s.mockService, s.logger)
	handler.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_CheckStatusError(t *testing.T) {
	s := &controllerTestSuite{}
	s.setup(t)
	defer s.teardown()

	expectedError := errors.New("status check error")

	s.mockService.EXPECT().CheckStatus(&s.request).Return(nil, expectedError)

	request, err := http.NewRequest("GET", "/status?orderId=order123&merchantOrderId=merchantOrder123", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := Handler(s.mockService, s.logger)
	handler.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
