package middlewares

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/authentication"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/logger"
)

//	{
//	 "login": "login",
//	 "id": 1,
//	 "role":{"id":2,"slug":"administrator"}
//	}
const testJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6ImxvZ2luIiwiaWQiOjEsInJvbGUiOnsiaWQiOjIsInNsdWciOiJhZG1pbmlzdHJhdG9yIn19.rdMalxI1tOIbyNeLAEmbSd4SYpSA42bcw6NswMn4iYo"

func TestBearerAuthenticateFunc(t *testing.T) {
	l, _ := logger.Get()
	f := BearerAuthenticateFunc("123", l)

	i, err := f(testJWT)
	assert.NoError(t, err)

	auth := i.(authentication.Auth)
	assert.Equal(t, 1, auth.Id)
	assert.Equal(t, "login", auth.Login)
	assert.Equal(t, 2, auth.Role.Id)
	assert.Equal(t, "administrator", auth.Role.Slug)
}
