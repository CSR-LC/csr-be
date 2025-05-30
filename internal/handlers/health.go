package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations/system"
)

func SetHealthHandler(logger *zap.Logger, api *operations.BeAPI) {
	petKindHandler := NewHealth(logger)

	api.SystemGetHealthHandler = petKindHandler.GetHealthFunc()

}

type Health struct {
	logger *zap.Logger
}

func NewHealth(logger *zap.Logger) *Health {
	return &Health{
		logger: logger,
	}
}

func (pk Health) GetHealthFunc() system.GetHealthHandlerFunc {
	return func(p system.GetHealthParams) middleware.Responder {
		return system.NewGetHealthOK()
	}
}
