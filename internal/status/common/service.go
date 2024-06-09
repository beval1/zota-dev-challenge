package common

import (
	"go.uber.org/zap"
	"zota-dev-challenge/internal/config"
	"zota-dev-challenge/internal/status/shared"
)

type ServiceInterface interface {
	CheckStatus(req *shared.ClientRequest) (*shared.Response, error)
}

type Service struct {
	logger       *zap.Logger
	config       *config.Config
	statusClient shared.StatusPaymentGateway
}

func NewService(logger *zap.Logger, config *config.Config, statusClient shared.StatusPaymentGateway) *Service {
	return &Service{logger: logger, config: config, statusClient: statusClient}
}

func (s *Service) CheckStatus(req *shared.ClientRequest) (*shared.Response, error) {
	serviceModel := shared.Request{
		ClientRequest: *req,
	}

	res, err := s.statusClient.CheckStatus(serviceModel)
	if err != nil {
		s.logger.Error("Failed to check status", zap.Error(err))
		return nil, err
	}
	return res, nil
}
