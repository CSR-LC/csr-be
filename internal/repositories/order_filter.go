package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/order"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/orderstatus"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/orderstatusname"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/user"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/middlewares"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/utils"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/pkg/domain"
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

func (r *orderFilterRepository) GetOrdersByActiveFilter(ctx context.Context, ownerId int,
	filter string) ([]*ent.Order, error) {
	fmt.Println("receiving orders by filter111111111111111111111 ")
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch filter {
	case "all":
		fmt.Println("receiving ALL orders.........")

		result, err := tx.Order.Query().
			Where(order.HasUsersWith(user.ID(ownerId))).
			WithUsers().WithOrderStatus().WithEquipments().
			All(ctx)
		if err != nil {
			return nil, err
		}

		return result, nil

	case "active":
		fmt.Println("receiving ACTIVE orders.........")
		// result, err := tx.Order.Query().
		// 	Where(order.HasUsersWith(user.ID(ownerId))).
		// 	Where(
		// 		order.Or(
		// 			order.HasOrderStatusWith(
		// 				orderstatus.
		// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusInReview)),
		// 			),
		// 			order.HasOrderStatusWith(
		// 				orderstatus.
		// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusInProgress)),
		// 			),
		// 			order.HasOrderStatusWith(
		// 				orderstatus.
		// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusApproved)),
		// 			),
		// 			order.HasOrderStatusWith(
		// 				orderstatus.
		// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusOverdue)),
		// 			),
		// 			order.HasOrderStatusWith(
		// 				orderstatus.
		// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusPrepared)),
		// 			),
		// 		),
		// 	).
		// 	WithUsers().WithOrderStatus().WithEquipments().
		// 	All(ctx)

		fmt.Println("ownerID", ownerId)
		result, err := tx.Order.Query().
			Where(order.HasUsersWith(user.ID(ownerId))).
			Where(
				order.HasOrderStatusWith(
					orderstatus.
						HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusInReview)),
				),
			).
			WithUsers().WithOrderStatus().WithEquipments().
			All(ctx)

		// result, err := tx.Order.Query().
		// 	Where(order.HasUsersWith(user.ID(ownerId))).
		// 	Where(
		// 		order.HasOrderStatusWith(
		// 			orderstatus.
		// 				HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusInReview)),
		// 		),
		// 	).Order(ent.Desc("current_date")).First(ctx).
		// 	WithUsers().WithOrderStatus().WithEquipments().
		// 	All(ctx)

		if err != nil {
			return nil, err
		}

		// var resultEND []*ent.Order
		// // owner := o.Edges.Users
		// for _, value := range result {
		// 	// var statusToOrder *models.OrderStatus
		// 	allStatuses := value.Edges.OrderStatus
		// 	fmt.Println("allStatusess", allStatuses)
		// 	fmt.Println("value.Edges.OrderStatus", value.Edges.OrderStatus)
		// 	if len(allStatuses) != 0 {
		// 		lastStatus := allStatuses[0]
		// 		fmt.Println("lastStatus", lastStatus)
		// 		if lastStatus.ID == 1 {
		// 			resultEND = append(resultEND, value)
		// 		}
		// 		// for _, s := range allStatuses {
		// 		// 	if lastStatus != domain.OrderStatusInReview {

		// 		// 	}
		// 		// 	if s.CurrentDate.After(lastStatus.CurrentDate) {
		// 		// 		lastStatus = s
		// 		// 	}
		// 		// }
		// 		// mappedStatus, err := MapStatus(id, lastStatus)
		// 		// if err != nil {
		// 		// 	log.Error("failed to map status", zap.Error(err))
		// 		// 	return nil, err
		// 		// }
		// 		// statusToOrder = mappedStatus
		// 	}
		// }

		// return resultEND, nil
		return result, nil

	case "completed":
		fmt.Println("receiving COMPLETED orders.........")

		result, err := tx.Order.Query().
			Where(order.HasUsersWith(user.ID(ownerId))).
			Where(
				order.Or(
					order.HasOrderStatusWith(
						orderstatus.
							HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusRejected)),
					),
					order.HasOrderStatusWith(
						orderstatus.
							HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusBlocked)),
					),
					order.HasOrderStatusWith(
						orderstatus.
							HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusClosed)),
					),
				),
			).
			WithUsers().WithOrderStatus().WithEquipments().
			All(ctx)

		if err != nil {
			return nil, err
		}

		return result, nil

	default:
		return nil, nil
	}

	// result, err := tx.Order.Query().
	// 	Where(order.HasUsersWith(user.ID(ownerId))).
	// 	Where(
	// 		order.Or(
	// 			// order.HasOrderStatusWith(
	// 			// 	orderstatus.
	// 			// 		HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusInReview)),
	// 			// ),
	// 			order.HasOrderStatusWith(
	// 				orderstatus.
	// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusRejected)),
	// 			),

	// 			order.HasOrderStatusWith(
	// 				orderstatus.
	// 					HasOrderStatusNameWith(orderstatusname.StatusEQ(domain.OrderStatusClosed)),
	// 			),
	// 		),
	// 	).
	// 	WithUsers().WithOrderStatus().WithEquipments().
	// 	All(ctx)

	// if err != nil {
	// 	return nil, err
	// }

	fmt.Println("status111")

	// result, err := tx.OrderStatus.Query().
	// Where(predicate.OrderStatus(orderstatusname.StatusIn(statuses...))).
	// All(ctx)

	return nil, nil

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
		WithOrderStatus().
		Order(orderFunc).Limit(limit).Offset(offset).
		WithUsers().WithOrderStatus().WithEquipments().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}
