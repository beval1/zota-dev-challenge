package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	ZotaMerchantId         string
	ZotaAPISecretKey       string
	ZotaEndpointId         string
	ZotaBaseUrl            string
	ZotaDepositCallBackUrl string
	ZotaDepositRedirectUrl string
	ENV                    string
}

func New(logger *zap.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", zap.Error(err))
	}

	env, err := godotenv.Read()
	if err != nil {
		logger.Error("Error reading .env file", zap.Error(err))
	}

	// Read environment variables and assign them to the configuration variables
	return &Config{
		ZotaMerchantId:         env["ZOTA_MERCHANT_ID"],
		ZotaAPISecretKey:       env["ZOTA_API_SECRET_KEY"],
		ZotaEndpointId:         env["ZOTA_ENDPOINT_ID"],
		ZotaBaseUrl:            env["ZOTA_BASE_URL"],
		ZotaDepositCallBackUrl: env["ZOTA_DEPOSIT_CALLBACK_URL"],
		ZotaDepositRedirectUrl: env["ZOTA_DEPOSIT_REDIRECT_URL"],
		ENV:                    env["ENVIRONMENT"],
	}
}
