package domain

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/CSR-LC/csr-be/internal/generated/ent"
)

type OrderOverdueCheckup interface {
	Checkup(ctx context.Context, cln *ent.Client, logger *zap.Logger)
	PeriodicalCheckup(ctx context.Context, overdueTimeCheckDuration time.Duration, cln *ent.Client, logger *zap.Logger)
}
