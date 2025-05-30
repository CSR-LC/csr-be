package repositories

import (
	"context"
	"time"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/registrationconfirm"
	"github.com/CSR-LC/csr-be/internal/generated/ent/user"
	"github.com/CSR-LC/csr-be/internal/middlewares"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

type registrationConfirmRepository struct {
}

func (rc *registrationConfirmRepository) CreateToken(ctx context.Context, token string, ttl time.Time, userID int) error {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return err
	}
	tokens, err := tx.RegistrationConfirm.Query().QueryUsers().Where(user.IDEQ(userID)).QueryRegistrationConfirm().All(ctx)
	if err != nil {
		return err
	}
	for _, t := range tokens {
		if errDelete := tx.RegistrationConfirm.DeleteOneID(t.ID).Exec(ctx); errDelete != nil {
			return errDelete
		}
	}
	_, err = tx.RegistrationConfirm.Create().SetToken(token).SetTTL(ttl).SetUsersID(userID).Save(ctx)
	return err
}

func (rc *registrationConfirmRepository) GetToken(ctx context.Context, token string) (*ent.RegistrationConfirm, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return tx.RegistrationConfirm.Query().Where(registrationconfirm.TokenEQ(token)).WithUsers().Only(ctx)
}

func (rc *registrationConfirmRepository) DeleteToken(ctx context.Context, token string) error {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return err
	}
	_, err = tx.RegistrationConfirm.Delete().Where(registrationconfirm.TokenEQ(token)).Exec(ctx)
	return err
}

func NewRegistrationConfirmRepository() domain.RegistrationConfirmRepository {
	return &registrationConfirmRepository{}
}
