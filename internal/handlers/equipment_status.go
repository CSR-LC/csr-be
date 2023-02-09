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

const (
	EQUIPMENT_UNDER_REPAIR_COMMENT_FOR_ORDER = "Equipment under repair"
)

func SetEquipmentStatusHandler(logger *zap.Logger, api *operations.BeAPI) {
	equipmentStatusRepo := repositories.NewEquipmentStatusRepository()
	orderStatusRepo := repositories.NewOrderStatusRepository()

	statusHandler := NewEquipmentStatus(logger)

	api.EquipmentStatusUpdateEquipmentStatusOnUnavailableHandler = statusHandler.PutEquipmentStatusInRepairFunc(equipmentStatusRepo, orderStatusRepo)
	api.EquipmentStatusUpdateEquipmentStatusOnAvailableHandler = statusHandler.PutEquipmentStatusRemoveFromRepairFunc(equipmentStatusRepo, orderStatusRepo)
	api.EquipmentStatusCheckEquipmentStatusHandler = statusHandler.GetEquipmentStatusCheckDatesFunc(equipmentStatusRepo)
	api.EquipmentStatusUpdateRepairedEquipmentStatusDatesHandler = statusHandler.PutEquipmentStatusEditDatesFunc(equipmentStatusRepo)
}

type EquipmentStatus struct {
	logger *zap.Logger
}

func NewEquipmentStatus(logger *zap.Logger) *EquipmentStatus {
	return &EquipmentStatus{
		logger: logger,
	}
}

func (c EquipmentStatus) GetEquipmentStatusCheckDatesFunc(
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
			ID:         &s.EquipmentstatusID,
		}

		orderResult, userResult, err := eqStatusRepository.GetOrderAndUserByDates(
			ctx, int(*data.ID), time.Time(*data.StartDate), time.Time(*data.EndDate))
		if err != nil {
			c.logger.Error("check equipment status by dates failed", zap.Error(err))
			return eqStatus.NewCheckEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't check equipment status by provided start date and end date"))
		}

		if orderResult == nil && userResult == nil {
			return eqStatus.NewCheckEquipmentStatusOK().WithPayload(
				&models.EquipmentStatusRepairConfirmationResponse{})
		}

		eqStatusResult, err := eqStatusRepository.GetEquipmentStatusByID(
			ctx, int(*data.ID))
		if err != nil {
			c.logger.Error("receiving equipment status by id failed", zap.Error(err))
			return eqStatus.NewCheckEquipmentStatusDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't find equipment status by provided id"))
		}

		orderID := int64(orderResult.ID)
		equipmentID := int64(eqStatusResult.Edges.Equipments.ID)
		return eqStatus.NewCheckEquipmentStatusOK().WithPayload(
			&models.EquipmentStatusRepairConfirmationResponse{
				Data: &models.EquipmentStatusRepairConfirmation{
					EquipmentStatusID: data.ID,
					EndDate:           (*strfmt.DateTime)(&eqStatusResult.EndDate),
					StartDate:         (*strfmt.DateTime)(&eqStatusResult.StartDate),
					StatusName:        &eqStatusResult.Edges.EquipmentStatusName.Name,
					OrderID:           &orderID,
					UserEmail:         &userResult.Email,
					EquipmentID:       &equipmentID,
				},
			})
	}
}

func (c EquipmentStatus) PutEquipmentStatusInRepairFunc(
	eqStatusRepository domain.EquipmentStatusRepository,
	orderStatusRepo domain.OrderStatusRepository) eqStatus.UpdateEquipmentStatusOnUnavailableHandlerFunc {
	return func(s eqStatus.UpdateEquipmentStatusOnUnavailableParams, access interface{}) middleware.Responder {
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
			ID:         &s.EquipmentstatusID,
		}

		updatedEqStatus, err := eqStatusRepository.Update(ctx, &data)
		if err != nil {
			c.logger.Error("update equipment status failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusOnUnavailableDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't update equipment status"))
		}

		orderResult, userResult, err := eqStatusRepository.GetOrderAndUserByDates(
			ctx, int(*data.ID), time.Time(*data.StartDate), time.Time(*data.EndDate))
		if err != nil {
			c.logger.Error("receiving user and order status by provided dates failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusOnUnavailableDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't receive order and user by provided start/end dates for updating equipment status"))
		}

		comment := EQUIPMENT_UNDER_REPAIR_COMMENT_FOR_ORDER
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
		return eqStatus.NewUpdateEquipmentStatusOnUnavailableOK().WithPayload(
			&models.EquipmentStatusRepairResponse{
				Data: &models.EquipmentStatusRepair{
					EquipmentStatusID: &id,
					EndDate:           (*strfmt.DateTime)(&updatedEqStatus.EndDate),
					StartDate:         (*strfmt.DateTime)(&updatedEqStatus.StartDate),
					StatusName:        newStatus,
				},
			})
	}
}

func (c EquipmentStatus) PutEquipmentStatusRemoveFromRepairFunc(
	eqStatusRepository domain.EquipmentStatusRepository,
	orderStatusRepo domain.OrderStatusRepository) eqStatus.UpdateEquipmentStatusOnAvailableHandlerFunc {
	return func(s eqStatus.UpdateEquipmentStatusOnAvailableParams, access interface{}) middleware.Responder {
		ctx := s.HTTPRequest.Context()
		newStatus := s.Name.StatusName

		ok, err := equipmentStatusAccessRights(access)
		if err != nil {
			c.logger.Error("Error while getting authorization", zap.Error(err))
			return orders.NewAddNewOrderStatusDefault(http.StatusInternalServerError).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Can't get authorization"}})
		}

		if !ok {
			c.logger.Error("User have no right to update equipment status on available", zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "You don't have rights to update equipment status"}})
		}

		if !newStatusIsAvailable(*newStatus) {
			c.logger.Error("Wrong new equipment status, status should be only 'available'", zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Wrong new equipment status, status should be only 'not available'"}})
		}

		timeNow := time.Now()
		data := models.EquipmentStatus{
			EndDate:    (*strfmt.DateTime)(&timeNow),
			StatusName: newStatus,
			ID:         &s.EquipmentstatusID,
		}

		updatedEqStatus, err := eqStatusRepository.Update(ctx, &data)
		if err != nil {
			c.logger.Error("update equipment on available status failed", zap.Error(err))
			return eqStatus.NewUpdateEquipmentStatusOnAvailableDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't update equipment status on available status"))
		}

		id := int64(updatedEqStatus.ID)
		return eqStatus.NewUpdateEquipmentStatusOnAvailableOK().WithPayload(
			&models.EquipmentStatusRepairResponse{
				Data: &models.EquipmentStatusRepair{
					EquipmentStatusID: &id,
					EndDate:           (*strfmt.DateTime)(&updatedEqStatus.EndDate),
					StartDate:         (*strfmt.DateTime)(&updatedEqStatus.StartDate),
					StatusName:        newStatus,
				},
			})
	}
}

func (c EquipmentStatus) PutEquipmentStatusEditDatesFunc(
	eqStatusRepository domain.EquipmentStatusRepository,
) eqStatus.UpdateRepairedEquipmentStatusDatesHandlerFunc {
	return func(s eqStatus.UpdateRepairedEquipmentStatusDatesParams, access interface{}) middleware.Responder {
		ctx := s.HTTPRequest.Context()

		ok, err := equipmentStatusAccessRights(access)
		if err != nil {
			c.logger.Error("Error while getting authorization", zap.Error(err))
			return orders.NewAddNewOrderStatusDefault(http.StatusInternalServerError).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "Can't get authorization"}})
		}

		if !ok {
			c.logger.Error("User have no right to update equipment status on available", zap.Any("access", access))
			return orders.NewAddNewOrderStatusDefault(http.StatusForbidden).
				WithPayload(&models.Error{Data: &models.ErrorData{Message: "You don't have rights to update equipment status"}})
		}

		reduceOneDayFromCurrentStartDate := strfmt.DateTime(
			time.Time(*s.Name.StartDate).AddDate(0, 0, -1),
		)

		addOneDayToCurrentEndDate := strfmt.DateTime(
			time.Time(*s.Name.EndDate).AddDate(0, 0, 1),
		)

		data := models.EquipmentStatus{
			StartDate: &reduceOneDayFromCurrentStartDate,
			EndDate:   &addOneDayToCurrentEndDate,
			ID:        &s.EquipmentstatusID,
		}

		updatedEqStatus, err := eqStatusRepository.Update(ctx, &data)
		if err != nil {
			c.logger.Error("update equipment on available status failed", zap.Error(err))
			return eqStatus.NewUpdateRepairedEquipmentStatusDatesDefault(http.StatusInternalServerError).
				WithPayload(buildStringPayload("can't update equipment status on available status"))
		}

		id := int64(updatedEqStatus.ID)
		return eqStatus.NewUpdateRepairedEquipmentStatusDatesOK().WithPayload(
			&models.EquipmentStatusRepairResponse{
				Data: &models.EquipmentStatusRepair{
					EquipmentStatusID: &id,
					EndDate:           (*strfmt.DateTime)(&updatedEqStatus.EndDate),
					StartDate:         (*strfmt.DateTime)(&updatedEqStatus.StartDate),
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

func newStatusIsAvailable(status string) bool {
	return status == domain.EquipmentStatusAvailable
}
