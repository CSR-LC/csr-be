package petkind

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/pet_kind"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests/common"
)

var (
	petKindName         = gofakeit.Name()
	migrationKindNumber = 3
)

func TestIntegration_GetPetKind(t *testing.T) {
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

	t.Run("get pet kind ok", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken

		petKindID, err := getKindIDByName(ctx, client, token, petKindName)
		require.NoError(t, err)

		paramsGet := pet_kind.NewGetPetKindParamsWithContext(ctx)
		paramsGet.SetPetKindID(*petKindID)

		got, err := client.PetKind.GetPetKind(paramsGet, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		want := pet_kind.NewGetPetKindOK()
		want.Payload = &models.PetKindResponse{
			ID:   petKindID,
			Name: &petKindName,
		}

		assert.Equal(t, got, want)
	})

	t.Run("Get Pet Kind By ID failed: incorrect ID", func(t *testing.T) {
		_, gotErr := client.PetKind.GetPetKind(pet_kind.NewGetPetKindParamsWithContext(ctx).WithPetKindID(-33), auth)
		require.Error(t, gotErr)

		wantErr := pet_kind.NewGetPetKindDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "Error while getting pet kind",
		}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Get Pet Kind By ID failed: access failed", func(t *testing.T) {
		token := utils.TokenNotExist
		_, gotErr := client.PetKind.GetPetKind(pet_kind.NewGetPetKindParamsWithContext(ctx).WithPetKindID(1), utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := pet_kind.NewGetPetKindDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_EditPetKind(t *testing.T) {
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

	token := loginUser.GetPayload().AccessToken

	petKindID, err := getKindIDByName(ctx, client, token, petKindName)
	require.NoError(t, err)

	auth := utils.AuthInfoFunc(loginUser.GetPayload().AccessToken)

	t.Run("Edit Kind By ID Ok: patch pet Kind parameters not empty, changed", func(t *testing.T) {
		params := pet_kind.NewEditPetKindParamsWithContext(ctx)
		petKind := "динозавр"

		params.PetKindID = *petKindID
		params.EditPetKind = &models.PetKind{
			Name: &petKind,
		}

		kind, err := client.PetKind.EditPetKind(params, auth)
		require.NoError(t, err)

		assert.Equal(t, petKind, *kind.GetPayload().Name)

		// revert changes for delete function
		params.EditPetKind = &models.PetKind{
			Name: &petKindName,
		}

		_, err = client.PetKind.EditPetKind(params, auth)
		assert.NoError(t, err)
	})

	t.Run("Edit Pet Kind By ID Ok: patch pet Kind parameters are empty, not changed", func(t *testing.T) {
		empty := ""
		patchKind := &models.PetKind{Name: &empty}
		kind, err := client.PetKind.EditPetKind(pet_kind.NewEditPetKindParamsWithContext(ctx).WithPetKindID(*petKindID).WithEditPetKind(patchKind), auth)
		require.NoError(t, err)
		assert.Equal(t, petKindName, *kind.GetPayload().Name)
	})

	t.Run("Edit Pet Kind By ID failed: ID incorrect", func(t *testing.T) {
		empty := ""
		patchKind := &models.PetKind{Name: &empty}
		_, gotErr := client.PetKind.EditPetKind(pet_kind.NewEditPetKindParamsWithContext(ctx).WithPetKindID(-33).WithEditPetKind(patchKind), auth)
		require.Error(t, gotErr)

		wantErr := pet_kind.NewEditPetKindDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "Error while updating pet kind",
		}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Edit Pet Kind By ID failed: access failed", func(t *testing.T) {
		empty := ""
		patchKind := &models.PetKind{Name: &empty}
		token := utils.TokenNotExist
		_, gotErr := client.PetKind.EditPetKind(pet_kind.NewEditPetKindParamsWithContext(ctx).WithPetKindID(*petKindID).WithEditPetKind(patchKind), utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := pet_kind.NewEditPetKindDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Edit Pet Kind By ID failed: validation failed, expect 422", func(t *testing.T) {
		_, gotErr := client.PetKind.EditPetKind(pet_kind.NewEditPetKindParamsWithContext(ctx).WithPetKindID(*petKindID), auth)
		require.Error(t, gotErr)

		wantErr := pet_kind.NewEditPetKindDefault(http.StatusUnprocessableEntity)
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

	t.Run("Delete Pet Kind By ID failed: not a validation error, pet kindID is not required in spec, expect 500", func(t *testing.T) {
		_, gotErr := client.PetKind.DeletePetKind(pet_kind.NewDeletePetKindParamsWithContext(ctx), auth)
		require.Error(t, gotErr)

		wantErr := pet_kind.NewDeletePetKindDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: "Error while deleting pet kind"}}
		assert.Equal(t, wantErr, gotErr)

		_, gotErr = client.PetKind.DeletePetKind(pet_kind.NewDeletePetKindParamsWithContext(ctx).WithPetKindID(-33), auth)
		require.Error(t, gotErr)

		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Delete Pet Kind By ID failed: access failed", func(t *testing.T) {
		token := utils.TokenNotExist
		_, gotErr := client.PetKind.DeletePetKind(pet_kind.NewDeletePetKindParamsWithContext(ctx).WithPetKindID(1), utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := pet_kind.NewDeletePetKindDefault(http.StatusInternalServerError)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Delete Pet Kind By ID OK", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken

		petKindID, err := getKindIDByName(ctx, client, token, petKindName)
		require.NoError(t, err)

		params := pet_kind.NewDeletePetKindParamsWithContext(ctx)
		params.PetKindID = *petKindID

		kind, err := client.PetKind.DeletePetKind(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		assert.Equal(t, kind.GetPayload(), "Pet kind deleted")
	})
}
