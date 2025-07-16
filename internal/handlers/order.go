package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/order"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/models"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations/orders"
	"github.com/CSR-LC/csr-be/internal/messages"
	"github.com/CSR-LC/csr-be/internal/repositories"
	"github.com/CSR-LC/csr-be/internal/utils"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

func SetOrderHandler(logger *zap.Logger, api *operations.BeAPI) {
	orderRepo := repositories.NewOrderRepository()
	eqStatusRepo := repositories.NewEquipmentStatusRepository()
	equipmentRepo := repositories.NewEquipmentRepository()
	ordersHandler := NewOrder(logger)

	api.OrdersGetUserOrdersHandler = ordersHandler.ListUserOrdersFunc(orderRepo)
	api.OrdersCreateOrderHandler = ordersHandler.CreateOrderFunc(orderRepo, eqStatusRepo, equipmentRepo)
	api.OrdersUpdateOrderHandler = ordersHandler.UpdateOrderFunc(orderRepo)
	api.OrdersGetAllOrdersHandler = ordersHandler.ListAllOrdersFunc(orderRepo)
	api.OrdersGetOrderHandler = ordersHandler.GetOrderFunc(orderRepo)
	api.OrdersDeleteOrderHandler = ordersHandler.DeleteOrderFunc(orderRepo)
}

type Order struct {
	logger *zap.Logger
}

func NewOrder(logger *zap.Logger) *Order {
	return &Order{
		logger: logger,
	}
}

func mapUserOrder(o *ent.Order, log *zap.Logger) (*models.UserOrder, error) {
	if o == nil {
		log.Warn("order is nil")
		return nil, errors.New("order is nil")
	}
	id := int64(o.ID)
	quantity := int64(o.Quantity)
	rentEnd := strfmt.DateTime(o.RentEnd)
	rentStart := strfmt.DateTime(o.RentStart)
	equipments := o.Edges.Equipments
	if equipments == nil {
		log.Warn("order has no equipments")
		return nil, errors.New("order has no equipments")
	}
	orderEquipments := make([]*models.EquipmentResponse, len(equipments))
	for i, eq := range equipments {
		var statusId int64
		var categoryId int64
		if eq.Edges.Category != nil {
			categoryId = int64(eq.Edges.Category.ID)
		}
		var subcategoryID int64
		if eq.Edges.Subcategory != nil {
			subcategoryID = int64(eq.Edges.Subcategory.ID)
		}
		if eq.Edges.CurrentStatus != nil {
			statusId = int64(eq.Edges.CurrentStatus.ID)
		}
		var photoID string
		if eq.Edges.Photo != nil {
			photoID = eq.Edges.Photo.ID
		}

		var psID int64
		eqID := int64(eq.ID)
		if eq.Edges.PetSize != nil {
			psID = int64(eq.Edges.PetSize.ID)
		}

		var petKinds []*models.PetKind
		if eq.Edges.PetKinds != nil {
			for _, petKind := range eq.Edges.PetKinds {
				j := &models.PetKind{
					Name: &petKind.Name,
				}
				petKinds = append(petKinds, j)
			}
		}

		var eqReceiptDate int64
		if eq.ReceiptDate != "" {
			eqReceiptTime, err := time.Parse(utils.TimeFormat, eq.ReceiptDate)
			if err != nil {
				log.Error("error during parsing date string")
				return nil, err
			}
			eqReceiptDate = eqReceiptTime.Unix()
		}

		orderEquipments[i] = &models.EquipmentResponse{
			TermsOfUse:       &eq.TermsOfUse,
			CompensationCost: &eq.CompensationCost,
			Condition:        eq.Condition,
			Description:      &eq.Description,
			ID:               &eqID,
			InventoryNumber:  &eq.InventoryNumber,
			Category:         &categoryId,
			Subcategory:      subcategoryID,
			Location:         nil,
			Name:             &eq.Name,
			PhotoID:          &photoID,
			PetSize:          &psID,
			PetKinds:         petKinds,
			ReceiptDate:      &eqReceiptDate,
			Supplier:         &eq.Supplier,
			TechnicalIssues:  &eq.TechIssue,
			Title:            &eq.Title,
			Status:           &statusId,
		}
	}

	var ownerId int64
	var ownerName string

	if o.Edges.Users != nil {
		ownerId = int64(o.Edges.Users.ID)
		ownerName = o.Edges.Users.Login
	} else {
		log.Warn("Order has no associated user", zap.Int("order_id", o.ID))
	}

	var statusToOrder *models.OrderStatus
	allStatuses := o.Edges.OrderStatus
	if len(allStatuses) != 0 {
		lastStatus := allStatuses[0]
		for _, s := range allStatuses {
			if s.CurrentDate.After(lastStatus.CurrentDate) {
				lastStatus = s
			}
		}
		mappedStatus, err := MapStatus(id, lastStatus)
		if err != nil {
			log.Error("failed to map status", zap.Error(err))
			return nil, err
		}
		statusToOrder = mappedStatus
	}

	return &models.UserOrder{
		Description: &o.Description,
		Equipments:  orderEquipments,
		ID:          &id,
		Quantity:    &quantity,
		RentEnd:     &rentEnd,
		RentStart:   &rentStart,
		User: &models.UserEmbeddable{
			ID:   &ownerId,
			Name: &ownerName,
		},
		LastStatus: statusToOrder,
		IsFirst:    &o.IsFirst,
	}, nil
}

func mapUserOrdersToResponse(log *zap.Logger, entOrders ...*ent.Order) ([]*models.UserOrder, error) {
	modelOrders := make([]*models.UserOrder, len(entOrders))
	for i, o := range entOrders {
		order, err := mapUserOrder(o, log)
		if err != nil {
			log.Error("failed to map order", zap.Error(err))
			return nil, err
		}
		modelOrders[i] = order
	}

	return modelOrders, nil
}

func mapOrdersToResponse(log *zap.Logger, entOrders ...*ent.Order) ([]*models.Order, error) {
	modelOrders := make([]*models.Order, len(entOrders))
	for i, o := range entOrders {
		uo, err := mapUserOrder(o, log)
		if err != nil {
			log.Error("failed to map order", zap.Error(err))
			return nil, err
		}
		user := mapUserInfoWoRole(o.Edges.Users)
		mo := &models.Order{
			Description: uo.Description,
			Equipments:  uo.Equipments,
			ID:          uo.ID,
			IsFirst:     uo.IsFirst,
			LastStatus:  uo.LastStatus,
			Quantity:    uo.Quantity,
			RentEnd:     uo.RentEnd,
			RentStart:   uo.RentStart,
			User:        user,
		}
		modelOrders[i] = mo
	}
	return modelOrders, nil
}

func (o Order) ListUserOrdersFunc(repository domain.OrderRepository) orders.GetUserOrdersHandlerFunc {
	return func(p orders.GetUserOrdersParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)
		limit := utils.GetValueByPointerOrDefaultValue(p.Limit, math.MaxInt)
		offset := utils.GetValueByPointerOrDefaultValue(p.Offset, 0)
		orderBy := utils.GetValueByPointerOrDefaultValue(p.OrderBy, utils.AscOrder)
		orderColumn := utils.GetValueByPointerOrDefaultValue(p.OrderColumn, order.FieldID)

		orderFilter := domain.OrderFilter{
			Filter: domain.Filter{
				Limit:       int(limit),
				Offset:      int(offset),
				OrderBy:     orderBy,
				OrderColumn: orderColumn,
			},
			Status: p.Status,
		}
		if p.Status != nil {
			_, ok := domain.AllOrderStatuses[*p.Status]
			if !ok {
				return orders.NewGetUserOrdersDefault(http.StatusBadRequest).
					WithPayload(buildBadRequestErrorPayload(messages.ErrQueryOrders, fmt.Sprintf("invalid order status '%v'", *p.Status)))
			}
		}

		total, err := repository.OrdersTotal(ctx, &userID)
		if err != nil {
			o.logger.Error(messages.ErrQueryTotalOrders, zap.Error(err))
			return orders.NewGetUserOrdersDefault(http.StatusInternalServerError).
				WithPayload(buildBadRequestErrorPayload(messages.ErrQueryTotalOrders, err.Error()))
		}

		var items []*ent.Order
		if total > 0 {
			items, err = repository.List(ctx, &userID, orderFilter)
			if err != nil {
				o.logger.Error(messages.ErrQueryOrders, zap.Error(err))
				return orders.NewGetUserOrdersDefault(http.StatusInternalServerError).
					WithPayload(buildBadRequestErrorPayload(messages.ErrQueryOrders, err.Error()))
			}
		}

		mappedOrders, err := mapUserOrdersToResponse(o.logger, items...)
		if err != nil {
			o.logger.Error(messages.ErrMapOrder, zap.Error(err))
			return orders.NewGetUserOrdersDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrMapOrder, err.Error()))
		}
		totalOrders := int64(total)
		listOrders := &models.UserOrdersList{
			Items: mappedOrders,
			Total: &totalOrders,
		}
		return orders.NewGetUserOrdersOK().WithPayload(listOrders)
	}
}

func (o Order) ListAllOrdersFunc(repository domain.OrderRepository) orders.GetAllOrdersHandlerFunc {
	return func(p orders.GetAllOrdersParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		limit := utils.GetValueByPointerOrDefaultValue(p.Limit, math.MaxInt)
		offset := utils.GetValueByPointerOrDefaultValue(p.Offset, 0)
		orderBy := utils.GetValueByPointerOrDefaultValue(p.OrderBy, utils.AscOrder)
		orderColumn := utils.GetValueByPointerOrDefaultValue(p.OrderColumn, order.FieldID)

		orderFilter := domain.OrderFilter{
			Filter: domain.Filter{
				Limit:       int(limit),
				Offset:      int(offset),
				OrderBy:     orderBy,
				OrderColumn: orderColumn,
			},
		}

		if p.Status != nil {
			_, ok := domain.AllOrderStatuses[*p.Status]
			if !ok {
				return orders.NewGetAllOrdersDefault(http.StatusBadRequest).
					WithPayload(buildBadRequestErrorPayload(messages.ErrQueryOrders, fmt.Sprintf("invalid order status '%v'", *p.Status)))
			}
			orderFilter.Status = p.Status
		}

		if p.EquipmentID != nil {
			eid := int(*p.EquipmentID)
			orderFilter.EquipmentID = &eid
		}

		total, err := repository.OrdersTotal(ctx, nil)
		if err != nil {
			o.logger.Error(messages.ErrQueryTotalOrders, zap.Error(err))
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrQueryTotalOrders, err.Error()))
		}

		var items []*ent.Order
		if total > 0 {
			items, err = repository.List(ctx, nil, orderFilter)
			if err != nil {
				o.logger.Error(messages.ErrQueryOrders, zap.Error(err))
				return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload(messages.ErrQueryOrders, err.Error()))
			}
		}

		mappedOrders, err := mapOrdersToResponse(o.logger, items...)
		if err != nil {
			o.logger.Error(messages.ErrMapOrder, zap.Error(err))
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrMapOrder, err.Error()))
		}
		totalOrders := int64(total)
		listOrders := &models.OrdersList{
			Items: mappedOrders,
			Total: &totalOrders,
		}
		return orders.NewGetAllOrdersOK().WithPayload(listOrders)
	}
}

func (o Order) CreateOrderFunc(
	orderRepo domain.OrderRepository,
	eqStatusRepo domain.EquipmentStatusRepository,
	equipmentRepo domain.EquipmentRepository,
) orders.CreateOrderHandlerFunc {
	return func(p orders.CreateOrderParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)

		id := int(*p.Data.EquipmentID)

		rentStart := time.Unix(0, *p.Data.RentStart)
		rentEnd := time.Unix(0, *p.Data.RentEnd)

		isEquipmentAvailable, err := eqStatusRepo.HasStatusByPeriod(ctx, domain.EquipmentStatusAvailable, id,
			rentStart, rentEnd)
		if err != nil {
			o.logger.Error(messages.ErrCheckEqStatusFailed, zap.Error(err))
			return orders.NewCreateOrderDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrCheckEqStatusFailed, err.Error()))
		}

		if !isEquipmentAvailable {
			o.logger.Warn(messages.ErrEquipmentIsNotFree)
			return orders.NewCreateOrderDefault(http.StatusConflict).
				WithPayload(buildConflictErrorPayload(messages.ErrEquipmentIsNotFree, ""))
		}

		if rentStart.After(rentEnd) {
			return orders.NewCreateOrderDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload(messages.ErrStartDateAfterEnd, ""))
		}

		if rentEnd.Sub(rentStart).Hours() < 24 {
			return orders.NewCreateOrderDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload(messages.ErrSmallRentPeriod, ""))
		}

		order, err := orderRepo.Create(ctx, p.Data, userID, []int{id})
		if err != nil {
			o.logger.Error(messages.ErrMapOrder, zap.Error(err))
			return orders.NewCreateOrderDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrMapOrder, err.Error()))
		}

		equipmentBookedStartDate := strfmt.DateTime(rentStart.AddDate(0, 0, -1))
		equipmentBookedEndDate := strfmt.DateTime(rentEnd.AddDate(0, 0, 1))
		eqID := int64(id)
		if _, err = eqStatusRepo.Create(ctx, &models.NewEquipmentStatus{
			EquipmentID: &eqID,
			StartDate:   &equipmentBookedStartDate,
			EndDate:     &equipmentBookedEndDate,
			StatusName:  &domain.EquipmentStatusBooked,
			OrderID:     int64(order.ID),
		}); err != nil {
			o.logger.Error(messages.ErrCreateEqStatus, zap.Error(err))
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrCreateEqStatus, err.Error()))
		}

		mappedOrder, err := mapUserOrder(order, o.logger)
		if err != nil {
			o.logger.Error(messages.ErrMapOrder, zap.Error(err))
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrMapOrder, err.Error()))
		}

		return orders.NewCreateOrderCreated().WithPayload(mappedOrder)
	}
}

func (o Order) UpdateOrderFunc(repository domain.OrderRepository) orders.UpdateOrderHandlerFunc {
	return func(p orders.UpdateOrderParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)
		orderID := int(p.OrderID)

		order, err := repository.Update(ctx, orderID, p.Data, userID)
		if err != nil {
			o.logger.Error(messages.ErrUpdateOrder, zap.Error(err))
			return orders.NewUpdateOrderDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrUpdateOrder, err.Error()))
		}

		mappedOrder, err := mapUserOrder(order, o.logger)
		if err != nil {
			o.logger.Error(messages.ErrMapOrder, zap.Error(err))
			return orders.NewUpdateOrderDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrMapOrder, err.Error()))
		}

		return orders.NewUpdateOrderOK().WithPayload(mappedOrder)
	}
}

func (o Order) GetOrderFunc(repository domain.OrderRepository) orders.GetOrderHandlerFunc {
	return func(p orders.GetOrderParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		orderID := int(p.OrderID)

		order, err := repository.Get(ctx, orderID)
		if err != nil {
			if ent.IsNotFound(err) {
				o.logger.Error(messages.ErrOrderNotFound, zap.Error(err))
				return orders.NewGetOrderNotFound().WithPayload(
					buildNotFoundErrorPayload(messages.ErrOrderNotFound, ""),
				)
			} else {
				o.logger.Error(messages.ErrGetOrder, zap.Error(err))
				return orders.NewGetOrderDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload(messages.ErrGetOrder, err.Error()))
			}
		}

		mappedOrders, err := mapOrdersToResponse(o.logger, order)
		if err != nil {
			o.logger.Error(messages.ErrGetOrder, zap.Error(err))
			return orders.NewGetOrderDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrGetOrder,
					fmt.Sprintf("Failed to map order to response: %v", err.Error())))
		}

		if len(mappedOrders) == 0 {
			o.logger.Error(messages.ErrOrderNotFound)
			return orders.NewGetOrderNotFound().WithPayload(
				buildNotFoundErrorPayload(messages.ErrOrderNotFound, ""),
			)
		}

		return orders.NewGetOrderOK().WithPayload(mappedOrders[0])
	}
}

func (o Order) DeleteOrderFunc(repository domain.OrderRepository) orders.DeleteOrderHandlerFunc {
	return func(p orders.DeleteOrderParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		orderID := int(p.OrderID)

		err := repository.Delete(ctx, orderID)
		if err != nil {
			if errors.Is(err, repositories.ErrOrderNotFound) {
				o.logger.Error(messages.ErrOrderNotFound, zap.Error(err))
				return orders.NewDeleteOrderNotFound().WithPayload(
					buildNotFoundErrorPayload(messages.ErrOrderNotFound, ""),
				)
			} else {
				o.logger.Error(messages.ErrDeleteOrder, zap.Error(err))
				return orders.NewDeleteOrderDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload(messages.ErrDeleteOrder, err.Error()))
			}
		}

		return orders.NewDeleteOrderNoContent()
	}
}
