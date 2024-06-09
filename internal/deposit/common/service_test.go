package common

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/deposit/common/zota"
	"zota-dev-challenge/internal/deposit/shared"
)

type serviceTestSuite struct {
	mockCtrl    *gomock.Controller
	mockGateway *zota.MockDepositPaymentGateway
	logger      *zap.Logger
	service     *Service
	request     shared.ClientRequest
}

func (s *serviceTestSuite) setup(t *testing.T) {
	s.mockCtrl = gomock.NewController(t)
	s.mockGateway = zota.NewMockDepositPaymentGateway(s.mockCtrl)
	s.logger, _ = zap.NewDevelopment()

	cfg := &config.Config{}

	s.service = NewService(s.logger, cfg, s.mockGateway)

	s.request = shared.ClientRequest{
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
	}
}

func (s *serviceTestSuite) teardown() {
	s.mockCtrl.Finish()
}

func TestProcessDeposit_Success(t *testing.T) {
	s := &serviceTestSuite{}
	s.setup(t)
	defer s.teardown()

	expectedResponse := &shared.Response{
		ClientRequest:         s.request,
		OrderID:               "order123",
		PaymentGatewayOrderID: "gateway123",
	}

	s.mockGateway.EXPECT().Deposit(shared.Request{
		ClientRequest: s.request,
	}).Return(expectedResponse, nil)

	response, err := s.service.ProcessDeposit(&s.request)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestProcessDeposit_InvalidCurrency(t *testing.T) {
	s := &serviceTestSuite{}
	s.setup(t)
	defer s.teardown()

	s.request.OrderCurrency = "EUR"

	response, err := s.service.ProcessDeposit(&s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid currency", err.Error())
}

func TestProcessDeposit_DepositError(t *testing.T) {
	s := &serviceTestSuite{}
	s.setup(t)
	defer s.teardown()

	expectedError := errors.New("deposit error")

	s.mockGateway.EXPECT().Deposit(shared.Request{
		ClientRequest: s.request,
	}).Return(nil, expectedError)

	response, err := s.service.ProcessDeposit(&s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, expectedError, err)
}
