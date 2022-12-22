package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetDB(t *testing.T) {
	_, _, err := GetDB(context.TODO(), "host=123", zap.NewNop())
	assert.ErrorContains(t, err, "failed to ping sql connection: failed to connect")
}
