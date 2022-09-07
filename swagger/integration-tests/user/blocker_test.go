package user

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/users"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests/common"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests/common"
)

var (
	auth runtime.ClientAuthInfoWriterFunc
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		ctx := context.Background()
		beClient := utils.SetupClient()

		var err error
		testLogin, testPassword, err = utils.GenerateLoginAndPassword()
		if err != nil {
			log.Fatalf("GenerateLoginAndPassword: %v", err)
		}
		user, err := utils.CreateUser(ctx, beClient, testLogin, testPassword)
		if err != nil {
			log.Fatalf("CreateUser: %v", err)
		}
		loginUser, err := utils.LoginUser(ctx, beClient, testLogin, testPassword)
		if err != nil {
			log.Fatalf("LoginUser: %v", err)
		}
		testUserID = *user.ID
		auth = utils.AuthInfoFunc(loginUser.GetPayload().AccessToken)

		os.Exit(m.Run())
	}
}

func TestIntegration_BlockUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	t.Run("Block User ok", func(t *testing.T) {
		params := users.NewBlockUserParamsWithContext(ctx)
		params.SetUserID(testUserID)

		_, err := client.Users.BlockUser(params, auth)
		require.NoError(t, err)

		res, err := client.Users.GetCurrentUser(users.NewGetCurrentUserParamsWithContext(ctx), auth)
		require.NoError(t, err)
		require.True(t, *res.Payload.IsBlocked)
	})

	t.Run("Block User failed: user ID not found", func(t *testing.T) {
		params := users.NewBlockUserParamsWithContext(ctx)
		params.SetUserID(-10)

		_, gotErr := client.Users.BlockUser(params, auth)
		require.Error(t, gotErr)

		wantErr := users.NewBlockUserDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: "failed to get user: ent: user not found"}}

		require.Equal(t, wantErr, gotErr)
	})

	t.Run("Block User failed: access failed", func(t *testing.T) {
		params := users.NewBlockUserParamsWithContext(ctx)
		params.SetUserID(testUserID)

		token := common.TokenNotExist
		_, gotErr := client.Users.BlockUser(params, common.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := users.NewBlockUserDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_UnblockUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	t.Run("Unblock User ok", func(t *testing.T) {
		params := users.NewUnblockUserParamsWithContext(ctx)
		params.SetUserID(testUserID)

		_, err := client.Users.UnblockUser(params, auth)
		require.NoError(t, err)

		res, err := client.Users.GetCurrentUser(users.NewGetCurrentUserParamsWithContext(ctx), auth)
		require.NoError(t, err)
		require.False(t, *res.Payload.IsBlocked)
	})

	t.Run("Unblock User failed: user ID not found", func(t *testing.T) {
		params := users.NewUnblockUserParamsWithContext(ctx)
		params.SetUserID(-10)

		_, gotErr := client.Users.UnblockUser(params, auth)
		require.Error(t, gotErr)

		wantErr := users.NewUnblockUserDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: "failed to get user: ent: user not found"}}

		require.Equal(t, wantErr, gotErr)
	})

	t.Run("Unblock User failed: access failed", func(t *testing.T) {
		params := users.NewUnblockUserParamsWithContext(ctx)
		params.SetUserID(testUserID)

		token := common.TokenNotExist
		_, gotErr := client.Users.UnblockUser(params, common.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := users.NewUnblockUserDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}
