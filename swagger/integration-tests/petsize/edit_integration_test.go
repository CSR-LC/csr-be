package petsize

import (
	"context"
	"net/http"
	"testing"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/pet_size"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_EditPetSize(t *testing.T) {
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
	token := loginUser.GetPayload().AccessToken

	petSizeID, err := getSizeIDByName(ctx, client, token, name)
	require.NoError(t, err)

	t.Run("Edit Pet Size By ID Ok: patch parameters not empty, changed", func(t *testing.T) {
		params := pet_size.NewEditPetSizeParamsWithContext(ctx)
		newSize := "new test size"
		newName := "new name"
		params.PetSizeID = *petSizeID
		params.EditPetSize = &models.PetSize{
			// should not provide ID
			Name: &newName,
			Size: &newSize,
		}

		ps, err := client.PetSize.EditPetSize(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		assert.Equal(t, newName, *ps.GetPayload().Name)
		assert.Equal(t, newSize, *ps.GetPayload().Size)

		// revert changes back for delete function
		params.EditPetSize = &models.PetSize{
			Name: &name,
			Size: &size,
		}

		_, err = client.PetSize.EditPetSize(params, utils.AuthInfoFunc(token))
		assert.NoError(t, err)
	})

	t.Run("Edit Pet Size By ID Ok: patch parameters are empty, not changed", func(t *testing.T) {
		list, err := client.PetSize.GetAllPetSize(pet_size.NewGetAllPetSizeParamsWithContext(ctx), auth)
		assert.NoError(t, err)

		id := list.Payload[0].ID
		empty := ""
		patch := &models.PetSize{
			Name: &empty,
			Size: &empty,
		}
		ps, err := client.PetSize.EditPetSize(pet_size.NewEditPetSizeParamsWithContext(ctx).WithPetSizeID(id).WithEditPetSize(patch), auth)
		require.NoError(t, err)
		assert.Equal(t, list.Payload[0].Name, ps.GetPayload().Name)
	})

	t.Run("Edit Pet Size By ID failed: ID incorrect", func(t *testing.T) {
		empty := ""
		patch := &models.PetSize{
			Name: &empty,
			Size: &empty,
		}
		_, gotErr := client.PetSize.EditPetSize(pet_size.NewEditPetSizeParamsWithContext(ctx).WithPetSizeID(-33).WithEditPetSize(patch), auth)
		require.Error(t, gotErr)

		wantErr := pet_size.NewEditPetSizeDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "Error while updating pet size",
		}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Edit Pet Size By ID failed: access failed", func(t *testing.T) {
		empty := ""
		patch := &models.PetSize{
			Name: &empty,
			Size: &empty,
		}
		token := utils.TokenNotExist
		_, gotErr := client.PetSize.EditPetSize(pet_size.NewEditPetSizeParamsWithContext(ctx).WithPetSizeID(1).WithEditPetSize(patch), utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := pet_size.NewEditPetSizeDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Edit Pet Size By ID failed: validation failed, expect 422", func(t *testing.T) {
		_, gotErr := client.PetSize.EditPetSize(pet_size.NewEditPetSizeParamsWithContext(ctx).WithPetSizeID(1), auth)
		require.Error(t, gotErr)

		wantErr := pet_size.NewEditPetSizeDefault(http.StatusUnprocessableEntity)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}
