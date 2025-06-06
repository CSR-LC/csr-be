package repositories

import (
	"context"
	"time"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
	"github.com/CSR-LC/csr-be/internal/generated/ent/emailconfirm"
	"github.com/CSR-LC/csr-be/internal/generated/ent/user"
	"github.com/CSR-LC/csr-be/internal/middlewares"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

type confirmEmailRepository struct {
}

func (p *confirmEmailRepository) CreateToken(ctx context.Context, token string, ttl time.Time, userID int, email string) error {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return err
	}

	tokens, err := tx.EmailConfirm.Query().QueryUsers().Where(user.IDEQ(userID)).
		QueryEmailConfirm().All(ctx)
	if err != nil {
		return err
	}

	for _, t := range tokens {
		if errDelete := tx.EmailConfirm.DeleteOneID(t.ID).Exec(ctx); errDelete != nil {
			return errDelete
		}
	}
	_, err = tx.EmailConfirm.Create().SetToken(token).SetTTL(ttl).SetEmail(email).
		SetUsersID(userID).Save(ctx)
	return err
}

func (p *confirmEmailRepository) GetToken(ctx context.Context, token string) (*ent.EmailConfirm, error) {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return tx.EmailConfirm.Query().Where(emailconfirm.TokenEQ(token)).WithUsers().Only(ctx)
}

func (p *confirmEmailRepository) DeleteToken(ctx context.Context, token string) error {
	tx, err := middlewares.TxFromContext(ctx)
	if err != nil {
		return err
	}
	_, err = tx.EmailConfirm.Delete().Where(emailconfirm.TokenEQ(token)).Exec(ctx)
	return err
}

func NewConfirmEmailRepository() domain.EmailConfirmRepository {
	return &confirmEmailRepository{}
}
