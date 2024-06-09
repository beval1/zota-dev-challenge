package common

import (
	"fmt"
	"go.uber.org/zap"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/deposit/shared"
)

type ServiceInterface interface {
	ProcessDeposit(r *shared.ClientRequest) (*shared.Response, error)
}

type Service struct {
	logger         *zap.Logger
	config         *config.Config
	depositGateway shared.DepositPaymentGateway // Use the interface
}

func NewService(logger *zap.Logger, config *config.Config, depositGateway shared.DepositPaymentGateway) *Service {
	return &Service{logger: logger, config: config, depositGateway: depositGateway}
}

func (s *Service) ProcessDeposit(req *shared.ClientRequest) (*shared.Response, error) {
	//validate the currency is USD as it is the only supported currency
	if req.OrderCurrency != "USD" {
		s.logger.Error("Invalid currency", zap.String("currency", req.OrderCurrency))
		return nil, fmt.Errorf("invalid currency")
	}

	//we use service model here to be easily extendable and decouple the service from the controller
	serviceModel := shared.Request{
		ClientRequest: *req,
	}

	depositRes, err := s.depositGateway.Deposit(serviceModel)
	if err != nil {
		s.logger.Error("Failed to process deposit", zap.Error(err))
		return nil, err
	}

	response := shared.Response{
		ClientRequest:         *req,
		OrderID:               depositRes.OrderID,
		PaymentGatewayOrderID: depositRes.PaymentGatewayOrderID,
	}
	return &response, nil
}
