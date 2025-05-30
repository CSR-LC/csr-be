package handlers

import (
	"math"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/order"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/models"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations/active_areas"
	"github.com/CSR-LC/csr-be/internal/messages"
	"github.com/CSR-LC/csr-be/internal/repositories"
	"github.com/CSR-LC/csr-be/internal/utils"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

func SetActiveAreaHandler(logger *zap.Logger, api *operations.BeAPI) {
	activeAreaRepo := repositories.NewActiveAreaRepository()
	activeAreaHandler := NewActiveArea(logger)
	api.ActiveAreasGetAllActiveAreasHandler = activeAreaHandler.GetActiveAreasFunc(activeAreaRepo)
}

type ActiveArea struct {
	logger *zap.Logger
}

func NewActiveArea(logger *zap.Logger) *ActiveArea {
	return &ActiveArea{
		logger: logger,
	}
}

func (area ActiveArea) GetActiveAreasFunc(repository domain.ActiveAreaRepository) active_areas.GetAllActiveAreasHandlerFunc {
	return func(a active_areas.GetAllActiveAreasParams, _ *models.Principal) middleware.Responder {
		ctx := a.HTTPRequest.Context()
		limit := utils.GetValueByPointerOrDefaultValue(a.Limit, math.MaxInt)
		offset := utils.GetValueByPointerOrDefaultValue(a.Offset, 0)
		orderBy := utils.GetValueByPointerOrDefaultValue(a.OrderBy, utils.AscOrder)
		orderColumn := utils.GetValueByPointerOrDefaultValue(a.OrderColumn, order.FieldID)
		total, err := repository.TotalActiveAreas(ctx)
		if err != nil {
			area.logger.Error(messages.ErrQueryTotalAreas, zap.Error(err))
			return active_areas.NewGetAllActiveAreasDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrQueryTotalAreas, err.Error()))
		}
		var e []*ent.ActiveArea
		if total > 0 {
			e, err = repository.AllActiveAreas(ctx, int(limit), int(offset), orderBy, orderColumn)
			if err != nil {
				area.logger.Error(messages.ErrQueryAreas, zap.Error(err))
				return active_areas.NewGetAllActiveAreasDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload(messages.ErrQueryAreas, err.Error()))
			}
		}
		totalAreas := int64(total)
		listActiveAreas := &models.ListOfActiveAreas{
			Items: make([]*models.ActiveArea, len(e)),
			Total: &totalAreas,
		}
		for i, element := range e {
			id := int64(element.ID)
			listActiveAreas.Items[i] = &models.ActiveArea{ID: &id, Name: &element.Name}
		}
		return active_areas.NewGetAllActiveAreasOK().WithPayload(listActiveAreas)
	}
}
