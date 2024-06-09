package common

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/status/common/zota"
	"zota-dev-challenge/internal/status/shared"
)

type statusTestSuite struct {
	mockCtrl    *gomock.Controller
	mockGateway *zota.MockStatusPaymentGateway
	logger      *zap.Logger
	service     *Service
	request     shared.ClientRequest
}

func (s *statusTestSuite) setup(t *testing.T) {
	s.mockCtrl = gomock.NewController(t)
	s.mockGateway = zota.NewMockStatusPaymentGateway(s.mockCtrl)
	s.logger, _ = zap.NewDevelopment()

	cfg := &config.Config{}

	s.service = NewService(s.logger, cfg, s.mockGateway)

	s.request = shared.ClientRequest{
		OrderId:         "1111",
		MerchantOrderId: "2222",
	}
}

func (s *statusTestSuite) teardown() {
	s.mockCtrl.Finish()
}

func TestCheckStatus_Success(t *testing.T) {
	s := &statusTestSuite{}
	s.setup(t)
	defer s.teardown()

	expectedResponse := &shared.Response{
		Status: "Success",
	}

	s.mockGateway.EXPECT().CheckStatus(shared.Request{
		ClientRequest: s.request,
	}).Return(expectedResponse, nil)

	response, err := s.service.CheckStatus(&s.request)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestCheckStatus_Error(t *testing.T) {
	s := &statusTestSuite{}
	s.setup(t)
	defer s.teardown()

	expectedError := errors.New("status check error")

	s.mockGateway.EXPECT().CheckStatus(shared.Request{
		ClientRequest: s.request,
	}).Return(nil, expectedError)

	response, err := s.service.CheckStatus(&s.request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, expectedError, err)
}
