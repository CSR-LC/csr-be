package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/order"
	"github.com/CSR-LC/csr-be/internal/generated/ent/orderstatus"
	"github.com/CSR-LC/csr-be/internal/generated/ent/orderstatusname"
	"github.com/CSR-LC/csr-be/internal/middlewares"
	"github.com/CSR-LC/csr-be/internal/utils"
)

type orderFilterRepository struct {
}

func NewOrderFilter() *orderFilterRepository {
	return &orderFilterRepository{}
}

var fieldsToOrderOrdersByStatus = []string{
	order.FieldID,
	order.FieldCreatedAt,
	order.FieldRentStart,
	order.FieldRentEnd,
}

func (r *orderFilterRepository) OrdersByStatusTotal(ctx context.Context, status string) (int, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return tx.OrderStatus.Query().
		QueryOrderStatusName().Where(orderstatusname.StatusEQ(status)).QueryOrderStatus().Count(ctx)
}

func (r *orderFilterRepository) OrdersByPeriodAndStatusTotal(ctx context.Context,
	from, to time.Time, status string) (int, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return tx.OrderStatus.Query().
		QueryOrderStatusName().Where(orderstatusname.StatusEQ(status)).QueryOrderStatus().
		Where(orderstatus.CurrentDateGT(from)).
		Where(orderstatus.CurrentDateLTE(to)).
		Count(ctx)
}

func (r *orderFilterRepository) OrdersByPeriodAndStatus(ctx context.Context, from, to time.Time, status string,
	limit, offset int, orderBy, orderColumn string) ([]*ent.Order, error) {
	if !utils.IsValueInList(orderColumn, fieldsToOrderOrdersByStatus) {
		return nil, errors.New("wrong field to order by")
	}
	orderFunc, err := utils.GetOrderFunc(orderBy, orderColumn)
	if err != nil {
		return nil, err
	}
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	items, err := tx.Order.Query().
		QueryOrderStatus().
		QueryOrderStatusName().Where(orderstatusname.StatusEQ(status)).
		QueryOrderStatus().
		Where(orderstatus.CurrentDateGT(from)).
		Where(orderstatus.CurrentDateLTE(to)).
		QueryOrder().
		WithOrderStatus().
		WithCurrentStatus().
		WithEquipments().
		WithUsers().
		Order(orderFunc).Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderFilterRepository) OrdersByStatus(ctx context.Context, status string,
	limit, offset int, orderBy, orderColumn string) ([]*ent.Order, error) {
	if !utils.IsValueInList(orderColumn, fieldsToOrderOrdersByStatus) {
		return nil, errors.New("wrong field to order by")
	}
	orderFunc, err := utils.GetOrderFunc(orderBy, orderColumn)
	if err != nil {
		return nil, err
	}
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	items, err := tx.Order.Query().
		QueryOrderStatus().
		QueryOrderStatusName().Where(orderstatusname.StatusEQ(status)).
		QueryOrderStatus().QueryOrder().
		Order(orderFunc).Limit(limit).Offset(offset).
		WithUsers().WithOrderStatus().WithEquipments().WithCurrentStatus().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}
