package petsize

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/pet_size"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests/common"
)

func TestIntegration_GetPetSize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	l, p, err := utils.GenerateLoginAndPassword()
	require.NoError(t, err)

	_, err = utils.CreateUser(ctx, client, l, p)
	require.NoError(t, err)

	loginUser, err := utils.LoginUser(ctx, client, l, p)
	require.NoError(t, err)

	auth := utils.AuthInfoFunc(loginUser.GetPayload().AccessToken)

	t.Run("get pet size ok", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken

		petSizeID, err := getSizeIDByName(ctx, client, token, name)
		require.NoError(t, err)

		paramsGet := pet_size.NewGetPetSizeParamsWithContext(ctx)
		paramsGet.SetPetSizeID(*petSizeID)

		got, err := client.PetSize.GetPetSize(paramsGet, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		want := pet_size.NewGetPetSizeOK()
		want.Payload = &models.PetSizeResponse{
			ID:   petSizeID,
			Name: &name,
			Size: &size,
		}

		assert.Equal(t, got, want)
	})

	t.Run("Get Pet Size OK: get universal, check isUniversal bool is true", func(t *testing.T) {
		list, err := client.PetSize.GetAllPetSize(pet_size.NewGetAllPetSizeParamsWithContext(ctx), auth)
		require.NoError(t, err)
		var id int64
		var psName *string

		for _, v := range list.Payload {
			if v.IsUniversal {
				id = v.ID
				psName = v.Name
			}
		}

		res, err := client.PetSize.GetPetSize(pet_size.NewGetPetSizeParamsWithContext(ctx).WithPetSizeID(id), auth)
		require.NoError(t, err)

		assert.True(t, res.Payload.IsUniversal)
		assert.Equal(t, psName, res.Payload.Name)
	})

	t.Run("Get Pet Size By ID failed: incorrect ID", func(t *testing.T) {
		_, gotErr := client.PetSize.GetPetSize(pet_size.NewGetPetSizeParamsWithContext(ctx).WithPetSizeID(-33), auth)
		require.Error(t, gotErr)

		wantErr := pet_size.NewGetPetSizeDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "Error while getting pet size",
		}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Get Pet Size By ID failed: access failed", func(t *testing.T) {
		token := utils.TokenNotExist
		_, gotErr := client.PetSize.GetPetSize(pet_size.NewGetPetSizeParamsWithContext(ctx).WithPetSizeID(1), utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := pet_size.NewGetPetSizeDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_DeletePetKind(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	l, p, err := utils.GenerateLoginAndPassword()
	require.NoError(t, err)

	_, err = utils.CreateUser(ctx, client, l, p)
	require.NoError(t, err)

	loginUser, err := utils.LoginUser(ctx, client, l, p)
	require.NoError(t, err)

	auth := utils.AuthInfoFunc(loginUser.GetPayload().AccessToken)

	t.Run("Delete Pet Size By ID failed: not a validation error, pet sizeID is not required in spec, expect 500", func(t *testing.T) {
		_, gotErr := client.PetSize.DeletePetSize(pet_size.NewDeletePetSizeParamsWithContext(ctx), auth)
		require.Error(t, gotErr)

		wantErr := pet_size.NewDeletePetSizeDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: "Error while deleting pet size"}}
		assert.Equal(t, wantErr, gotErr)

		_, gotErr = client.PetSize.DeletePetSize(pet_size.NewDeletePetSizeParamsWithContext(ctx).WithPetSizeID(-33), auth)
		require.Error(t, gotErr)

		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Delete Pet Size By ID failed: access failed", func(t *testing.T) {
		token := utils.TokenNotExist
		_, gotErr := client.PetSize.DeletePetSize(pet_size.NewDeletePetSizeParamsWithContext(ctx).WithPetSizeID(1), utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := pet_size.NewDeletePetSizeDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Delete Pet Size By ID OK", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken

		petSizeID, err := getSizeIDByName(ctx, client, token, name)
		require.NoError(t, err)

		params := pet_size.NewDeletePetSizeParamsWithContext(ctx)
		params.PetSizeID = *petSizeID

		res, err := client.PetSize.DeletePetSize(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		assert.Equal(t, "Pet size deleted", res.GetPayload())
	})
}

func getSizeIDByName(ctx context.Context, client *client.Be, token *string, petSizeName string) (*int64, error) {
	paramsGetAll := pet_size.NewGetAllPetSizeParamsWithContext(ctx)
	allPetSize, err := client.PetSize.GetAllPetSize(paramsGetAll, utils.AuthInfoFunc(token))
	if err != nil {
		return nil, err
	}
	var petSizeID *int64

	for _, petSize := range allPetSize.GetPayload() {
		if *petSize.Name == petSizeName {
			petSizeID = &petSize.ID
		}
	}
	return petSizeID, nil
}
