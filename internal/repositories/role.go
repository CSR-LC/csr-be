package repositories

import (
	"context"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/role"
	"github.com/CSR-LC/csr-be/internal/middlewares"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

type roleRepository struct {
}

func (r *roleRepository) GetRoles(ctx context.Context) ([]*ent.Role, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return tx.Role.Query().Order(ent.Asc(role.FieldID)).All(ctx)
}

func NewRoleRepository() domain.RoleRepository {
	return &roleRepository{}
}
