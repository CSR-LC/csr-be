package orders

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/orders"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests"
)

func TestIntegration_CreateOrder(t *testing.T) {
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

	t.Run("Create Order", func(t *testing.T) {
		params := orders.NewCreateOrderParamsWithContext(ctx)
		desc := "test description"
		quantity := int64(1)
		equipment := int64(1)
		rentStart := strfmt.DateTime(time.Now())
		rentEnd := strfmt.DateTime(time.Now().Add(time.Hour * 24))
		params.Data = &models.OrderCreateRequest{
			Equipment:   &equipment,
			Description: &desc,
			Quantity:    &quantity,
			RentStart:   &rentStart,
			RentEnd:     &rentEnd,
		}
		// в базе нет значений для создания ордера
		res, err := client.Orders.CreateOrder(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		_ = res
	})
}
