package statuses

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/status"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testLogin    string
	testPassword string
)
var name = "test status name"

func TestMain(m *testing.M) {
	ctx := context.Background()
	beClient := utils.SetupClient()

	var err error
	testLogin, testPassword, err = utils.GenerateLoginAndPassword()
	if err != nil {
		log.Fatalf("GenerateLoginAndPassword: %v", err)
	}
	_, err = utils.CreateUser(ctx, beClient, testLogin, testPassword)
	if err != nil {
		log.Fatalf("CreateUser: %v", err)
	}
	os.Exit(m.Run())
}

func TestIntegration_PostStatus(t *testing.T) {
	ctx := context.Background()
	beClient := utils.SetupClient()

	loginUser, err := utils.LoginUser(ctx, beClient, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("post status ok", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken

		params := status.NewPostStatusParamsWithContext(ctx)
		params.Name = &models.StatusName{
			Name: &name,
		}
		res, err := beClient.Status.PostStatus(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		assert.Equal(t, name, res.GetPayload().Data.Name)
	})
}

func TestIntegration_GetAllStatuses(t *testing.T) {
	ctx := context.Background()
	beClient := utils.SetupClient()

	loginUser, err := utils.LoginUser(ctx, beClient, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("get all statuses ok", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken

		params := status.NewGetStatusesParams()
		_, err = beClient.Status.GetStatuses(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)
	})
}

func TestIntegration_GetStatus(t *testing.T) {
	ctx := context.Background()
	beClient := utils.SetupClient()

	loginUser, err := utils.LoginUser(ctx, beClient, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("get status by id ok", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken
		id, err := getStatusIDByName(ctx, beClient, token, name)
		require.NoError(t, err)
		require.NotNil(t, id)

		params := status.NewGetStatusParamsWithContext(ctx)
		params.StatusID = *id
		res, err := beClient.Status.GetStatus(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		assert.Equal(t, *id, res.GetPayload().Data.ID)
		assert.Equal(t, &name, res.GetPayload().Data.Name)
	})
}

func TestIntegration_Delete(t *testing.T) {
	ctx := context.Background()
	beClient := utils.SetupClient()

	loginUser, err := utils.LoginUser(ctx, beClient, testLogin, testPassword)
	require.NoError(t, err)

	token := loginUser.GetPayload().AccessToken
	t.Run("delete status by id passed", func(t *testing.T) {
		id, err := getStatusIDByName(ctx, beClient, token, fmt.Sprintf("%s-updated", name))
		require.NoError(t, err)
		require.NotNil(t, id)

		params := status.NewDeleteStatusParamsWithContext(ctx)
		params.StatusID = *id

		_, err = beClient.Status.DeleteStatus(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)
	})
}

func getStatusIDByName(ctx context.Context, client *client.Be, token *string, name string) (*int64, error) {
	paramsGetAll := status.NewGetStatusesParamsWithContext(ctx)
	all, err := client.Status.GetStatuses(paramsGetAll, utils.AuthInfoFunc(token))
	if err != nil {
		return nil, err
	}
	var id *int64

	for _, st := range all.GetPayload() {
		if *st.Name == name {
			id = st.ID
		}
	}
	return id, nil
}
