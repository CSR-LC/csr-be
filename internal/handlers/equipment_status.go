package handlers

import (
	"net/http"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/authentication"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations"
	eqStatus "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations/equipment_status"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations/orders"
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
		newStatus := s.Name.StatusName

		ok, err := equipmentStatusAccessRights(access)
		if err != nil {
			c.logger.Error("Error while getting authorization", zap.Error(err))
			return orders.NewAddNewOrderStatusDefault(http.StatusInternalServerError).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Can't get authorization"}})
		}

		if !ok {
			c.logger.Error("User have no right to update equipment status", zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "You don't have rights to update equipment status"}})
		}

		if !checkStatus(*newStatus) {
			c.logger.Error("Wrong new equipment status, status should be only 'not available'", zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Wrong new equipment status, status should be only 'not available'"}})
		}

		data := models.EquipmentStatus{
			EndDate:    s.Name.EndDate,
			StartDate:  s.Name.StartDate,
			StatusName: newStatus,
			ID:         &s.ID,
		}

		updatedEqStatus, err := eqStatusRepository.Update(ctx, &data)
		if err != nil {
			c.logger.Error("update equipment status failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't update equipment status"))
		}

		id := int64(updatedEqStatus.ID)
		endTime := strfmt.DateTime(updatedEqStatus.EndDate)
		startTime := strfmt.DateTime(updatedEqStatus.StartDate)

		return eqStatus.NewUpdateEquipmentStatusOK().WithPayload(
			&models.EquipmentStatusUpdateResponse{
				Data: &models.EquipmentStatusUpdate{
					ID:         &id,
					EndDate:    &endTime,
					StartDate:  &startTime,
					StatusName: newStatus,
				},
			})
	}
}

func equipmentStatusAccessRights(access interface{}) (bool, error) {
	isManager, err := authentication.IsManager(access)
	if err != nil {
		return false, err
	}

	return isManager, nil
}

func checkStatus(status string) bool {
	return status == domain.EquipmentStatusNotAvailable
}
