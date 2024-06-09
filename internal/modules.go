package internal

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"zota-dev-challenge/internal/config"
	deposit "zota-dev-challenge/internal/deposit/common"
	zotaDeposit "zota-dev-challenge/internal/deposit/common/zota"
	depositShared "zota-dev-challenge/internal/deposit/shared"
	status "zota-dev-challenge/internal/status/common"
	zotaStatus "zota-dev-challenge/internal/status/common/zota"
	statusShared "zota-dev-challenge/internal/status/shared"
)

var AppModules = fx.Options(
	fx.Provide(func(logger *zap.Logger, config *config.Config) statusShared.StatusPaymentGateway {
		return zotaStatus.NewStatusGateway(logger, config)
	}),
	fx.Provide(func(logger *zap.Logger, config *config.Config) depositShared.DepositPaymentGateway {
		return zotaDeposit.NewDepositGateway(logger, config)
	}),
	fx.Provide(status.NewService),
	fx.Provide(deposit.NewService),
	fx.Provide(config.New),
	fx.Provide(validator.New),
	fx.Provide(InitRouterV1),
	fx.Provide(InitLogger),
)
