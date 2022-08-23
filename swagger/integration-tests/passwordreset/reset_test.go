package passwordreset

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/password_reset"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests"
)

func TestIntegration_Reset(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	l, p, err := utils.GenerateLoginAndPassword()
	require.NoError(t, err)

	_, err = utils.CreateUser(ctx, client, l, p)
	require.NoError(t, err)

	t.Run("Send Link By Login successfully", func(t *testing.T) {
		params := password_reset.NewSendLinkByLoginParamsWithContext(ctx)
		params.Login = &models.SendPasswordResetLinkRequest{
			Data: &models.Login{Login: &l},
		}
		got, err := client.PasswordReset.SendLinkByLogin(params)
		require.NoError(t, err)

		want := &password_reset.SendLinkByLoginOK{
			Payload: models.PasswordResetResponse("Reset link sent"),
		}
		assert.Equal(t, want, got)
	})

	t.Run("Get Password Reset Link failed: wrong token", func(t *testing.T) {
		params := password_reset.NewGetPasswordResetLinkParamsWithContext(ctx)
		params.Token = "some-dummy-token"
		_, err = client.PasswordReset.GetPasswordResetLink(params)
		require.Error(t, err)
		assert.Contains(t, "Failed to verify token. Please try again later", err.Error())
	})
}
