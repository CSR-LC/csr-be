package db

import (
	"testing"

	"github.com/CSR-LC/csr-be/internal/config"

	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {
	_, _, err := GetDB(config.DB{Host: "localhost"})
	require.ErrorContains(t, err, "failed to ping sql connection: failed to connect")
}
