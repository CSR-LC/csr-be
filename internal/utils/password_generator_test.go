package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPasswordGenerator_OK(t *testing.T) {
	generator, err := NewPasswordGenerator(10)
	assert.NoError(t, err)
	assert.NotNil(t, generator)
}

func TestPasswordGenerator_Generate(t *testing.T) {
	length := 10
	generator, err := NewPasswordGenerator(length)
	assert.NoError(t, err)
	password, err := generator.NewPassword()
	assert.NoError(t, err)
	assert.Len(t, password, length)
}
