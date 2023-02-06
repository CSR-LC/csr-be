package handlers

import (
	"net/http"
	"time"

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
	orderStatusRepo := repositories.NewOrderStatusRepository()

	statusHandler := NewEquipmentStatus(logger)

	api.EquipmentStatusUpdateEquipmentStatusHandler = statusHandler.PutEquipmentStatusFunc(equipmentStatusRepo, orderStatusRepo)
	api.EquipmentStatusCheckEquipmentStatusHandler = statusHandler.GetEquipmentStatusFunc(equipmentStatusRepo)
}

type EquipmentStatus struct {
	logger *zap.Logger
}

func NewEquipmentStatus(logger *zap.Logger) *EquipmentStatus {
	return &EquipmentStatus{
		logger: logger,
	}
}

func (c EquipmentStatus) GetEquipmentStatusFunc(
	eqStatusRepository domain.EquipmentStatusRepository) eqStatus.CheckEquipmentStatusHandlerFunc {
	return func(s eqStatus.CheckEquipmentStatusParams, access interface{}) middleware.Responder {
		ctx := s.HTTPRequest.Context()
		newStatus := s.Name.StatusName

		ok, err := equipmentStatusAccessRights(access)
		if err != nil {
			c.logger.Error("Error while getting authorization", zap.Error(err))
			return orders.NewAddNewOrderStatusDefault(http.StatusInternalServerError).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Can't get authorization"}})
		}

		if !ok {
			c.logger.Error("User have no right to check that equipment status has orders for provided dates",
				zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{
					Message: "You don't have rights to update equipment status"}},
				)
		}

		if !newStatusIsUnavailable(*newStatus) {
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

		statusResult, orderResult, userResult, err := eqStatusRepository.GetOrderAndUserByDates(
			ctx, int(*data.ID), time.Time(*data.StartDate), time.Time(*data.EndDate))
		if err != nil {
			c.logger.Error("check equipment status by dates failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't check equipment status by provided start date and end date"))
		}

		if orderResult == nil && userResult == nil {
			return eqStatus.NewCheckEquipmentStatusOK().WithPayload(
				&models.EquipmentStatusUpdateConfirmationResponse{})
		}

		orderID := int64(orderResult.ID)
		return eqStatus.NewCheckEquipmentStatusOK().WithPayload(
			&models.EquipmentStatusUpdateConfirmationResponse{
				Data: &models.EquipmentStatusUpdateConfirmation{
					EquipmentStatusID: data.ID,
					EndDate:           (*strfmt.DateTime)(&statusResult.EndDate),
					StartDate:         (*strfmt.DateTime)(&statusResult.StartDate),
					StatusName:        &statusResult.Edges.EquipmentStatusName.Name,
					OrderID:           &orderID,
					UserEmail:         &userResult.Email,
				},
			})
	}
}

func (c EquipmentStatus) PutEquipmentStatusFunc(
	eqStatusRepository domain.EquipmentStatusRepository,
	orderStatusRepo domain.OrderStatusRepository) eqStatus.UpdateEquipmentStatusHandlerFunc {
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

		if !newStatusIsUnavailable(*newStatus) {
			c.logger.Error("Wrong new equipment status, status should be only 'not available'", zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Wrong new equipment status, status should be only 'not available'"}})
		}

		reduceOneDayFromCurrentStartDate := strfmt.DateTime(
			time.Time(*s.Name.StartDate).AddDate(0, 0, -1),
		)

		addOneDayToCurrentEndDate := strfmt.DateTime(
			time.Time(*s.Name.EndDate).AddDate(0, 0, 1),
		)

		data := models.EquipmentStatus{
			StartDate:  &reduceOneDayFromCurrentStartDate,
			EndDate:    &addOneDayToCurrentEndDate,
			StatusName: newStatus,
			ID:         &s.ID,
		}

		updatedEqStatus, err := eqStatusRepository.Update(ctx, &data)
		if err != nil {
			c.logger.Error("update equipment status failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't update equipment status"))
		}

		_, orderResult, userResult, err := eqStatusRepository.GetOrderAndUserByDates(
			ctx, int(*data.ID), time.Time(*data.StartDate), time.Time(*data.EndDate))
		if err != nil {
			c.logger.Error("receiving user and order status by provided dates failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't receive order and user by provided start/end dates for updating equipment status"))
		}

		comment := "equipment in repair"
		timeNow := time.Now()
		orderID := int64(orderResult.ID)
		model := models.NewOrderStatus{
			Comment:   &comment,
			CreatedAt: (*strfmt.DateTime)(&timeNow),
			OrderID:   &orderID,
			Status:    &domain.OrderStatusRejected,
		}

		err = orderStatusRepo.UpdateStatus(ctx, userResult.ID, model)
		if err != nil {
			c.logger.Error("Update order status error", zap.Error(err))
			return orders.NewAddNewOrderStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("Can't update order status"))
		}

		id := int64(updatedEqStatus.ID)
		return eqStatus.NewUpdateEquipmentStatusOK().WithPayload(
			&models.EquipmentStatusUpdateResponse{
				Data: &models.EquipmentStatusUpdate{
					ID:         &id,
					EndDate:    (*strfmt.DateTime)(&updatedEqStatus.EndDate),
					StartDate:  (*strfmt.DateTime)(&updatedEqStatus.StartDate),
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

func newStatusIsUnavailable(status string) bool {
	return status == domain.EquipmentStatusNotAvailable
}

// func newStatusIsAvailable(status string) bool {
// 	return status == domain.EquipmentStatusAvailable
// }
