package repositories

import (
	"context"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/equipmentstatusname"
	"github.com/CSR-LC/csr-be/internal/middlewares"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

type equipmentStatusNameRepository struct{}

func NewEquipmentStatusNameRepository() domain.EquipmentStatusNameRepository {
	return &equipmentStatusNameRepository{}
}

func (r *equipmentStatusNameRepository) Create(ctx context.Context, name string) (*ent.EquipmentStatusName, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return tx.EquipmentStatusName.Create().SetName(name).Save(ctx)
}

func (r *equipmentStatusNameRepository) GetAll(ctx context.Context) ([]*ent.EquipmentStatusName, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return tx.EquipmentStatusName.Query().All(ctx)
}

func (r *equipmentStatusNameRepository) Get(ctx context.Context, id int) (*ent.EquipmentStatusName, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return tx.EquipmentStatusName.Get(ctx, id)
}
func (r *equipmentStatusNameRepository) GetByName(ctx context.Context, name string) (*ent.EquipmentStatusName, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return tx.EquipmentStatusName.Query().Where(equipmentstatusname.NameEQ(name)).First(ctx)
}

func (r *equipmentStatusNameRepository) Delete(ctx context.Context, id int) (*ent.EquipmentStatusName, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	statusToDelete, err := tx.EquipmentStatusName.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.EquipmentStatusName.DeleteOne(statusToDelete).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return statusToDelete, nil
}
