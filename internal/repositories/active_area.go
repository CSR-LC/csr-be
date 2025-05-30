package repositories

import (
	"context"
	"errors"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/activearea"
	"github.com/CSR-LC/csr-be/internal/middlewares"
	"github.com/CSR-LC/csr-be/internal/utils"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

var fieldsToOrderAreas = []string{
	activearea.FieldID,
	activearea.FieldName,
}

type activeAreaRepository struct {
}

func NewActiveAreaRepository() domain.ActiveAreaRepository {
	return &activeAreaRepository{}
}

func (r *activeAreaRepository) AllActiveAreas(ctx context.Context, limit, offset int,
	orderBy, orderColumn string) ([]*ent.ActiveArea, error) {
	if !utils.IsValueInList(orderColumn, fieldsToOrderAreas) {
		return nil, errors.New("wrong column to order by")
	}
	orderFunc, err := utils.GetOrderFunc(orderBy, orderColumn)
	if err != nil {
		return nil, err
	}
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return tx.ActiveArea.Query().Order(orderFunc).Limit(limit).Offset(offset).All(ctx)
}

func (r *activeAreaRepository) TotalActiveAreas(ctx context.Context) (int, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return tx.ActiveArea.Query().Count(ctx)
}
