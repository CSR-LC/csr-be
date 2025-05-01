package repositories

import (
	"context"
	"fmt"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/middlewares"
)

type orderStatusNameRepository struct {
}

func NewOrderStatusNameRepository() *orderStatusNameRepository {
	return &orderStatusNameRepository{}
}

func (r *orderStatusNameRepository) ListOfOrderStatusNames(ctx context.Context) ([]*ent.OrderStatusName, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	pointersStatuses, err := tx.OrderStatusName.Query().Order(ent.Asc("id")).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("status history error, failed to get status names: %s", err)
	}
	return pointersStatuses, nil
}
