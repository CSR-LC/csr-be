package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateFileName() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	name := strings.ReplaceAll(id.String(), "-", "")
	return name, nil
}
