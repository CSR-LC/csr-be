package handlers

import (
	"net/http"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations"
	eqStatus "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations/equipment_status"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/repositories"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/pkg/domain"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
)

func SetEquipmentStatusHandler(logger *zap.Logger, api *operations.BeAPI) {
	equipmentStatusRepo := repositories.NewEquipmentStatusRepository()
	statusHandler := NewEquipmentStatus(logger)

	api.EquipmentStatusUpdateEquipmentStatusHandler = statusHandler.PutEquipmentStatusFunc(equipmentStatusRepo)
}

type EquipmentStatus struct {
	logger *zap.Logger
}

func NewEquipmentStatus(logger *zap.Logger) *EquipmentStatus {
	return &EquipmentStatus{
		logger: logger,
	}
}

func (c EquipmentStatus) PutEquipmentStatusFunc(eqStatusRepository domain.EquipmentStatusRepository) eqStatus.UpdateEquipmentStatusHandlerFunc {
	return func(s eqStatus.UpdateEquipmentStatusParams, access interface{}) middleware.Responder {
		ctx := s.HTTPRequest.Context()
		statusName := s.Name.StatusName
		data := models.EquipmentStatus{}
		data.EndDate = s.Name.EndDate
		data.StartDate = s.Name.StartDate
		data.StatusName = statusName
		data.ID = &s.ID

		updatedEqStatus, err := eqStatusRepository.Update(ctx, &data)
		if err != nil {
			c.logger.Error("create status failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't create status"))
		}

		id := int64(updatedEqStatus.ID)
		endTime := strfmt.DateTime(updatedEqStatus.EndDate)
		startTime := strfmt.DateTime(updatedEqStatus.StartDate)
		return eqStatus.NewUpdateEquipmentStatusOK().WithPayload(&models.EquipmentStatusUpdateResponse{
			Data: &models.EquipmentStatusUpdate{
				ID:         &id,
				EndDate:    &endTime,
				StartDate:  &startTime,
				StatusName: statusName,
			},
		})
	}
}
