package handlers

import (
	"context"
	"errors"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/authentication"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations/orders"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/repositories"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
	"net/http"
)

type Order struct {
	client *ent.Client
	logger *zap.Logger
}

func NewOrder(client *ent.Client, logger *zap.Logger) *Order {
	return &Order{
		client: client,
		logger: logger,
	}
}

func mapOrder(ctx context.Context, o *ent.Order) (*models.Order, error) {
	id := int64(o.ID)
	quantity := int64(o.Quantity)
	rentEnd := strfmt.DateTime(o.RentEnd)
	rentStart := strfmt.DateTime(o.RentStart)
	owner, err := o.QueryUsers().Only(ctx)
	if err != nil {
		return nil, err
	}
	equipment, err := o.QueryEquipments().Only(ctx)
	if err != nil {
		return nil, err
	}
	ownerId := int64(owner.ID)
	ownerName := owner.Name
	var kindId int64
	if equipment.Edges.Kind != nil {
		kindId = int64(equipment.Edges.Kind.ID)
	}
	var statusId int64
	if equipment.Edges.Status != nil {
		statusId = int64(equipment.Edges.Status.ID)
	}

	return &models.Order{
		Description: &o.Description,
		Equipment: &models.Equipment{
			Description: &equipment.Description,
			Kind:        &kindId,
			Location:    nil,
			Name:        &equipment.Name,
			Photo:       nil,
			RateDay:     &equipment.RateDay,
			RateHour:    &equipment.RateHour,
			Sku:         &equipment.Sku,
			Status:      &statusId,
		},
		ID:        &id,
		Quantity:  &quantity,
		RentEnd:   &rentEnd,
		RentStart: &rentStart,
		User: &models.UserEmbeddable{
			ID:   &ownerId,
			Name: &ownerName,
		},
	}, nil
}

func mapOrdersToResponse(ctx context.Context, entOrders []*ent.Order) ([]*models.Order, error) {
	modelOrders := make([]*models.Order, len(entOrders))
	for i, o := range entOrders {
		order, err := mapOrder(ctx, o)
		if err != nil {
			return nil, err
		}
		modelOrders[i] = order
	}

	return modelOrders, nil
}

func (o Order) ListOrderFunc(repository repositories.OrderRepository) orders.GetAllOrdersHandlerFunc {
	return func(p orders.GetAllOrdersParams, access interface{}) middleware.Responder {
		ownerId, err := authentication.GetUserId(access)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, repositories.OrderAccessDenied{}) {
				status = http.StatusForbidden
			}
			return orders.NewGetAllOrdersDefault(status).WithPayload(buildErrorPayload(err))
		}
		ctx := p.HTTPRequest.Context()
		items, total, err := repository.List(ctx, ownerId)
		if err != nil {
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		mappedOrders, err := mapOrdersToResponse(ctx, items)
		if err != nil {
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		return orders.NewGetAllOrdersOK().WithPayload(&orders.GetAllOrdersOKBody{Data: &orders.GetAllOrdersOKBodyData{
			Items: mappedOrders,
			Total: int64(total),
		}})
	}
}

func (o Order) CreateOrderFunc(repository repositories.OrderRepository) orders.CreateOrderHandlerFunc {
	return func(p orders.CreateOrderParams, access interface{}) middleware.Responder {
		ownerId, err := authentication.GetUserId(access)
		if err != nil {
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}
		ctx := p.HTTPRequest.Context()
		order, err := repository.Create(ctx, p.Data, ownerId)
		if err != nil {
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		mappedOrder, err := mapOrder(ctx, order)
		if err != nil {
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		return orders.NewCreateOrderCreated().WithPayload(mappedOrder)
	}
}

func (o Order) UpdateOrderFunc(repository repositories.OrderRepository) orders.UpdateOrderHandlerFunc {
	return func(p orders.UpdateOrderParams, access interface{}) middleware.Responder {
		ownerId, err := authentication.GetUserId(access)
		if err != nil {
			return orders.NewGetAllOrdersDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}
		ctx := p.HTTPRequest.Context()
		id := int(p.OrderID)
		order, err := repository.Update(ctx, id, p.Data, ownerId)
		if err != nil {
			return orders.NewUpdateOrderDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		mappedOrder, err := mapOrder(ctx, order)
		if err != nil {
			return orders.NewUpdateOrderDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		return orders.NewUpdateOrderOK().WithPayload(mappedOrder)
	}
}
