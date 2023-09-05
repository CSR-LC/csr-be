package equipment

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/categories"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/equipment"
	eqStatusName "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/equipment_status_name"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/orders"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/pet_kind"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/pet_size"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/photos"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/client/subcategories"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/handlers"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/integration-tests/common"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/pkg/domain"
)

func TestIntegration_CreateEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)
	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)

	t.Run("Create Equipment", func(t *testing.T) {
		params := equipment.NewCreateNewEquipmentParamsWithContext(ctx)
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)

		params.NewEquipment = model

		res, err := client.Equipment.CreateNewEquipment(params, auth)
		require.NoError(t, err)

		// location returned as nil, in discussion we decided that this parameter will have more that one values
		// for now it is not handled
		// todo: uncomment string below when it's handled properly
		assert.Equal(t, model.Category, res.Payload.Category)
		assert.Equal(t, model.CompensationCost, res.Payload.CompensationCost)
		assert.Equal(t, model.Condition, res.Payload.Condition)
		assert.Equal(t, model.Description, res.Payload.Description)
		assert.Equal(t, model.InventoryNumber, res.Payload.InventoryNumber)
		assert.Equal(t, model.Category, res.Payload.Category)
		//assert.Equal(t, location, *res.Payload.Location)
		assert.Equal(t, model.MaximumDays, res.Payload.MaximumDays)
		assert.Equal(t, model.Name, res.Payload.Name)
		assert.Equal(t, model.PetSize, res.Payload.PetSize)
		assert.Contains(t, *res.Payload.PhotoID, *model.PhotoID)
		assert.Equal(t, model.ReceiptDate, res.Payload.ReceiptDate)
		assert.Equal(t, model.Status, res.Payload.Status)
		assert.Equal(t, model.Supplier, res.Payload.Supplier)
		assert.Equal(t, model.TechnicalIssues, res.Payload.TechnicalIssues)
		assert.Equal(t, model.Title, res.Payload.Title)
	})

	t.Run("Create Equipment failed: 422 status code error, description and name fields have a number of characters greater than the limit ",
		func(t *testing.T) {
			params := equipment.NewCreateNewEquipmentParamsWithContext(ctx)
			model, err := setParameters(ctx, client, auth)
			require.NoError(t, err)

			// name field tests:
			// max length of name field: 100 characters
			name, err := utils.GenerateRandomString(101)
			require.NoError(t, err)
			model.Name = &name
			params.NewEquipment = model

			_, err = client.Equipment.CreateNewEquipment(params, auth)
			require.Error(t, err)

			name, err = utils.GenerateRandomString(99)
			require.NoError(t, err)
			model.Name = &name
			params.NewEquipment = model

			_, err = client.Equipment.CreateNewEquipment(params, auth)
			require.NoError(t, err)

			// description field tests:
			// max length of description field: 255 characters
			model, err = setParameters(ctx, client, auth)
			require.NoError(t, err)
			description, err := utils.GenerateRandomString(256)
			require.NoError(t, err)
			model.Description = &description

			params.NewEquipment = model
			_, err = client.Equipment.CreateNewEquipment(params, auth)
			require.Error(t, err)

			description, err = utils.GenerateRandomString(254)
			require.NoError(t, err)
			model.Description = &description
			params.NewEquipment = model

			_, err = client.Equipment.CreateNewEquipment(params, auth)
			require.NoError(t, err)
		})

	t.Run("Create Equipment failed: foreign key constraint error", func(t *testing.T) {
		params := equipment.NewCreateNewEquipmentParamsWithContext(ctx)
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)

		id := ""
		model.PhotoID = &id
		params.NewEquipment = model

		_, gotErr := client.Equipment.CreateNewEquipment(params, auth)
		require.Error(t, gotErr)

		wantErr := equipment.NewCreateNewEquipmentDefault(500)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "Error while creating equipment",
		}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Create Equipment failed: authorization error 500 Invalid token", func(t *testing.T) {
		params := equipment.NewCreateNewEquipmentParamsWithContext(ctx)
		token := utils.TokenNotExist
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)

		params.NewEquipment = model

		_, gotErr := client.Equipment.CreateNewEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := equipment.NewCreateNewEquipmentDefault(http.StatusUnauthorized)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_GetAllEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)

	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)

	t.Run("Get All Equipment", func(t *testing.T) {
		params := equipment.NewGetAllEquipmentParamsWithContext(ctx)

		res, err := client.Equipment.GetAllEquipment(params, auth)
		require.NoError(t, err)
		assert.NotZero(t, len(res.Payload.Items))
	})

	t.Run("Get All Equipment: authorization error 500 Invalid token", func(t *testing.T) {
		params := equipment.NewGetAllEquipmentParamsWithContext(ctx)
		token := utils.TokenNotExist

		_, gotErr := client.Equipment.GetAllEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := equipment.NewGetAllEquipmentDefault(http.StatusUnauthorized)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_GetEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)

	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)

	model, err := setParameters(ctx, client, auth)
	require.NoError(t, err)

	created, err := createEquipment(ctx, client, auth, model)
	require.NoError(t, err)

	t.Run("Get Equipment", func(t *testing.T) {
		params := equipment.NewGetEquipmentParamsWithContext(ctx)
		require.NoError(t, err)

		params.EquipmentID = *created.Payload.ID
		res, err := client.Equipment.GetEquipment(params, auth)
		require.NoError(t, err)

		// location returned as nil, in discussion we decided that this parameter will have more that one values
		// for now it is not handled
		// todo: uncomment string below when it's handled properly
		assert.Equal(t, model.Category, res.Payload.Category)
		assert.Equal(t, model.CompensationCost, res.Payload.CompensationCost)
		assert.Equal(t, model.Condition, res.Payload.Condition)
		assert.Equal(t, model.Description, res.Payload.Description)
		assert.Equal(t, model.InventoryNumber, res.Payload.InventoryNumber)
		assert.Equal(t, model.Category, res.Payload.Category)
		assert.Equal(t, model.MaximumDays, res.Payload.MaximumDays)
		assert.Equal(t, model.Name, res.Payload.Name)
		//assert.Equal(t, model.Location, res.Payload.Location)
		assert.Equal(t, model.PetSize, res.Payload.PetSize)
		assert.Contains(t, *res.Payload.PhotoID, *model.PhotoID)
		assert.Equal(t, model.ReceiptDate, res.Payload.ReceiptDate)
		assert.Equal(t, model.Status, res.Payload.Status)
		assert.Equal(t, model.Supplier, res.Payload.Supplier)
		assert.Equal(t, model.TechnicalIssues, res.Payload.TechnicalIssues)
		assert.Equal(t, model.Title, res.Payload.Title)
	})

	t.Run("Get Equipment failed: passed invalid equipment id", func(t *testing.T) {
		params := equipment.NewGetEquipmentParamsWithContext(ctx)
		require.NoError(t, err)

		params.EquipmentID = int64(-10)
		_, gotErr := client.Equipment.GetEquipment(params, auth)
		require.Error(t, gotErr)

		wantErr := equipment.NewGetEquipmentDefault(500)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: "Error while getting equipment"}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Get Equipment failed: authorization error 500 Invalid token", func(t *testing.T) {
		params := equipment.NewGetEquipmentParamsWithContext(ctx)
		require.NoError(t, err)

		params.EquipmentID = *created.Payload.ID
		token := utils.TokenNotExist
		_, gotErr := client.Equipment.GetEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := equipment.NewGetEquipmentDefault(http.StatusUnauthorized)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_FindEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)

	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)

	model, err := setParameters(ctx, client, auth)
	require.NoError(t, err)

	_, err = createEquipment(ctx, client, auth, model)
	require.NoError(t, err)

	t.Run("Find Equipment", func(t *testing.T) {
		params := equipment.NewFindEquipmentParamsWithContext(ctx)
		params.FindEquipment = &models.EquipmentFilter{
			Category: *model.Category,
		}
		res, err := client.Equipment.FindEquipment(params, auth)
		require.NoError(t, err)

		assert.NotZero(t, *res.Payload.Total)
		for _, item := range res.Payload.Items {
			assert.Equal(t, model.Category, item.Category)
		}
	})

	t.Run("Find Equipment: limit = 1", func(t *testing.T) {
		params := equipment.NewFindEquipmentParamsWithContext(ctx)
		params.FindEquipment = &models.EquipmentFilter{
			Category: *model.Category,
		}
		limit := int64(1)
		params.WithLimit(&limit)

		res, err := client.Equipment.FindEquipment(params, auth)
		require.NoError(t, err)

		assert.Equal(t, int(limit), len(res.Payload.Items))
	})

	t.Run("Find Equipment failed: authorization error 500 Invalid token", func(t *testing.T) {
		params := equipment.NewFindEquipmentParamsWithContext(ctx)
		params.FindEquipment = &models.EquipmentFilter{
			Category: *model.Category,
		}

		token := utils.TokenNotExist
		_, gotErr := client.Equipment.FindEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := equipment.NewFindEquipmentDefault(http.StatusUnauthorized)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Find Equipment: unknown parameters, zero items found", func(t *testing.T) {
		params := equipment.NewFindEquipmentParamsWithContext(ctx)
		params.FindEquipment = &models.EquipmentFilter{
			TermsOfUse: "unknown category",
		}

		res, gotErr := client.Equipment.FindEquipment(params, auth)
		require.NoError(t, gotErr)

		assert.Zero(t, len(res.Payload.Items))
	})
}

func TestIntegration_EditEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)

	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
	model, err := setParameters(ctx, client, auth)
	require.NoError(t, err)

	created, err := createEquipment(ctx, client, auth, model)
	require.NoError(t, err)

	t.Run("Edit Equipment description", func(t *testing.T) {
		desc := "new description"
		model.Description = &desc
		params := equipment.NewEditEquipmentParamsWithContext(ctx).WithEquipmentID(*created.Payload.ID).
			WithEditEquipment(model)

		res, err := client.Equipment.EditEquipment(params, auth)
		require.NoError(t, err)

		assert.Equal(t, desc, *res.Payload.Description)
	})

	t.Run("Edit Equipment description failed: wrong Equipment ID", func(t *testing.T) {
		desc := "new description"
		model.Description = &desc
		params := equipment.NewEditEquipmentParamsWithContext(ctx).WithEquipmentID(int64(-10)).
			WithEditEquipment(model)

		_, gotErr := client.Equipment.EditEquipment(params, auth)
		require.Error(t, gotErr)

		wantErr := equipment.NewEditEquipmentDefault(500)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: "Error while updating equipment"}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Edit Equipment failed: authorization error 500 Invalid token", func(t *testing.T) {
		params := equipment.NewEditEquipmentParamsWithContext(ctx)
		require.NoError(t, err)

		params.EquipmentID = *created.Payload.ID
		token := utils.TokenNotExist
		_, gotErr := client.Equipment.EditEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := equipment.NewEditEquipmentDefault(http.StatusUnauthorized)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_ArchiveEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)

	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
	model, err := setParameters(ctx, client, auth)
	require.NoError(t, err)

	created, err := createEquipment(ctx, client, auth, model)
	require.NoError(t, err)

	t.Run("Archive Equipment", func(t *testing.T) {
		params := equipment.NewArchiveEquipmentParamsWithContext(ctx).WithEquipmentID(*created.Payload.ID)
		res, gotError := client.Equipment.ArchiveEquipment(params, auth)
		require.NoError(t, gotError)

		require.True(t, res.IsCode(http.StatusNoContent))
	})

	t.Run("Archive Equipment failed: equipment not found", func(t *testing.T) {
		params := equipment.NewArchiveEquipmentParamsWithContext(ctx).WithEquipmentID(-1)
		resp, gotErr := client.Equipment.ArchiveEquipment(params, auth)
		require.Error(t, gotErr)
		fmt.Print(resp)

		wantedErr := equipment.NewArchiveEquipmentNotFound()
		wantedErr.Payload = &models.Error{Data: &models.ErrorData{Message: handlers.EquipmentNotFoundMsg}}

		require.Equal(t, wantedErr, gotErr)
	})

	t.Run("Archive Equipment with active orders", func(t *testing.T) {
		var orderID *int64
		orderID, err = createOrder(ctx, client, auth, created.Payload.ID)
		params := equipment.NewArchiveEquipmentParamsWithContext(ctx).WithEquipmentID(*created.Payload.ID)
		var res *equipment.ArchiveEquipmentNoContent
		res, err = client.Equipment.ArchiveEquipment(params, auth)
		require.NoError(t, err)
		require.True(t, res.IsCode(http.StatusNoContent))
		var ok bool
		ok, err = checkOrderStatus(ctx, client, auth, orderID, "closed")
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("Archive Equipment failed: auth failed", func(t *testing.T) {
		params := equipment.NewArchiveEquipmentParamsWithContext(ctx).WithEquipmentID(*created.Payload.ID)
		token := utils.TokenNotExist

		_, gotErr := client.Equipment.ArchiveEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantedErr := equipment.NewArchiveEquipmentDefault(http.StatusUnauthorized)
		wantedErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantedErr, gotErr)
	})

	// todo: test for archive equipment with non-default status
}

func TestIntegration_DeleteEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()

	tokens := utils.AdminUserLogin(t)

	auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)

	t.Run("Delete All Equipment", func(t *testing.T) {
		res, err := client.Equipment.GetAllEquipment(equipment.NewGetAllEquipmentParamsWithContext(ctx), auth)
		require.NoError(t, err)
		assert.NotZero(t, len(res.Payload.Items))

		params := equipment.NewDeleteEquipmentParamsWithContext(ctx)
		for _, item := range res.Payload.Items {
			params.WithEquipmentID(*item.ID)
			_, err = client.Equipment.DeleteEquipment(params, auth)
			require.NoError(t, err)
		}

		res, err = client.Equipment.GetAllEquipment(equipment.NewGetAllEquipmentParamsWithContext(ctx), auth)
		require.NoError(t, err)
		assert.Zero(t, len(res.Payload.Items))
	})

	t.Run("Delete Equipment failed: zero equipments, delete failed", func(t *testing.T) {
		params := equipment.NewDeleteEquipmentParamsWithContext(ctx)
		params.WithEquipmentID(int64(1))
		_, gotErr := client.Equipment.DeleteEquipment(params, auth)
		require.Error(t, gotErr)

		wantErr := equipment.NewDeleteEquipmentDefault(500)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "Error while getting equipment",
		}}
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("Delete Equipment failed: auth failed", func(t *testing.T) {
		params := equipment.NewDeleteEquipmentParamsWithContext(ctx)
		params.WithEquipmentID(int64(1))
		token := utils.TokenNotExist

		_, gotErr := client.Equipment.DeleteEquipment(params, utils.AuthInfoFunc(&token))
		require.Error(t, gotErr)

		wantErr := equipment.NewDeleteEquipmentDefault(http.StatusUnauthorized)
		wantErr.Payload = &models.Error{Data: nil}
		assert.Equal(t, wantErr, gotErr)
	})
}

func TestIntegration_BlockEquipment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := utils.SetupClient()
	startDate := strfmt.DateTime(time.Now())
	endDate := strfmt.DateTime(time.Now().AddDate(0, 0, 10))

	t.Run("Block Equipment is prohibited for operators", func(t *testing.T) {
		tokens := utils.OperatorUserLogin(t)
		auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)
		eq, err := createEquipment(ctx, client, auth, model)
		require.NoError(t, err)

		params := equipment.NewBlockEquipmentParamsWithContext(ctx).WithEquipmentID(*eq.Payload.ID)
		params.Data = &models.ChangeEquipmentStatusToBlockedRequest{
			StartDate: strfmt.DateTime(startDate),
			EndDate:   strfmt.DateTime(endDate),
		}

		_, err = client.Equipment.BlockEquipment(params, auth)
		require.Error(t, err)

		wantErr := equipment.NewBlockEquipmentDefault(http.StatusForbidden)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "You don't have rights to block the equipment",
		}}
		assert.Equal(t, wantErr, err)
	})

	t.Run("Block Equipment is prohibited for admins", func(t *testing.T) {
		tokens := utils.AdminUserLogin(t)
		auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)
		eq, err := createEquipment(ctx, client, auth, model)
		require.NoError(t, err)

		params := equipment.NewBlockEquipmentParamsWithContext(ctx).WithEquipmentID(*eq.Payload.ID)
		params.Data = &models.ChangeEquipmentStatusToBlockedRequest{
			StartDate: strfmt.DateTime(startDate),
			EndDate:   strfmt.DateTime(endDate),
		}

		_, err = client.Equipment.BlockEquipment(params, auth)
		require.Error(t, err)

		wantErr := equipment.NewBlockEquipmentDefault(http.StatusForbidden)
		wantErr.Payload = &models.Error{Data: &models.ErrorData{
			Message: "You don't have rights to block the equipment",
		}}
		assert.Equal(t, wantErr, err)
	})

	t.Run("Block Equipment is permitted for managers", func(t *testing.T) {
		tokens := utils.ManagerUserLogin(t)
		auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)
		eq, err := createEquipment(ctx, client, auth, model)
		require.NoError(t, err)

		params := equipment.NewBlockEquipmentParamsWithContext(ctx).WithEquipmentID(*eq.Payload.ID)
		params.Data = &models.ChangeEquipmentStatusToBlockedRequest{
			StartDate: strfmt.DateTime(startDate),
			EndDate:   strfmt.DateTime(endDate),
		}

		res, err := client.Equipment.BlockEquipment(params, auth)
		require.NoError(t, err)
		require.True(t, res.IsCode(http.StatusNoContent))
	})

	t.Run("Block Equipment with active orders", func(t *testing.T) {
		tokens := utils.ManagerUserLogin(t)
		auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
		model, err := setParameters(ctx, client, auth)
		require.NoError(t, err)
		eq, err := createEquipment(ctx, client, auth, model)
		require.NoError(t, err)

		orderID1, err := createOrder(ctx, client, auth, eq.Payload.ID)
		require.NoError(t, err)
		fmt.Println(err)
		createOrder(ctx, client, auth, eq.Payload.ID)
		require.NoError(t, err)
		//fmt.Println(orderID2)
		listParams := orders.NewGetAllOrdersParamsWithContext(ctx)
		test, err := client.Orders.GetAllOrders(listParams, auth)
		require.NoError(t, err)
		fmt.Println("12122", test)
		for i, o := range test.Payload.Items {
			var st string
			if i == 0 {
				st = domain.OrderStatusApproved
			} else {
				st = domain.OrderStatusRejected
			}

			dt := strfmt.DateTime(time.Now())
			osp := orders.NewAddNewOrderStatusParamsWithContext(ctx)
			osp.Data = &models.NewOrderStatus{
				OrderID:   o.ID,
				CreatedAt: &dt,
				Status:    &st,
				Comment:   &st,
			}
			_, err = client.Orders.AddNewOrderStatus(osp, auth)
			require.NoError(t, err)
		}

		params := equipment.NewBlockEquipmentParamsWithContext(ctx).WithEquipmentID(*eq.Payload.ID)
		params.Data = &models.ChangeEquipmentStatusToBlockedRequest{
			StartDate: strfmt.DateTime(startDate),
			EndDate:   strfmt.DateTime(endDate),
		}

		var res *equipment.BlockEquipmentNoContent
		res, err = client.Equipment.BlockEquipment(params, auth)
		require.NoError(t, err)
		require.True(t, res.IsCode(http.StatusNoContent))

		orders, err := client.Orders.GetOrdersByStatus(
			orders.NewGetOrdersByStatusParamsWithContext(ctx).WithStatus(domain.OrderStatusBlocked), auth)
		fmt.Println(orders)

		ok, err := checkOrderStatus(ctx, client, auth, orderID1, domain.OrderStatusBlocked)
		fmt.Println(ok, err)
		//ok, err = checkOrderStatus(ctx, client, auth, orderID2, domain.OrderStatusBlocked)
		//fmt.Println(ok, err)
		//require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("Block Equipment is failed, equipment not found", func(t *testing.T) {
		tokens := utils.ManagerUserLogin(t)
		auth := utils.AuthInfoFunc(tokens.GetPayload().AccessToken)
		var fakeID int64 = 111

		params := equipment.NewBlockEquipmentParamsWithContext(ctx).WithEquipmentID(fakeID)
		params.Data = &models.ChangeEquipmentStatusToBlockedRequest{
			StartDate: strfmt.DateTime(startDate),
			EndDate:   strfmt.DateTime(endDate),
		}

		_, err := client.Equipment.BlockEquipment(params, auth)
		require.Error(t, err)

		wantErr := equipment.NewBlockEquipmentNotFound()
		wantErr.Payload = &models.Error{Data: &models.ErrorData{Message: handlers.EquipmentNotFoundMsg}}
		assert.Equal(t, wantErr, err)
	})
}

func createOrder(ctx context.Context, be *client.Be, auth runtime.ClientAuthInfoWriterFunc, id *int64) (*int64, error) {
	rentStart := strfmt.NewDateTime()
	dateTimeFmt := "2006-01-02 15:04:05"
	err := rentStart.UnmarshalText([]byte(time.Now().Format(dateTimeFmt)))
	if err != nil {
		return nil, err
	}
	rentEnd := strfmt.NewDateTime()
	err = rentEnd.UnmarshalText([]byte(time.Now().AddDate(0, 0, 10).Format(dateTimeFmt)))
	if err != nil {
		return nil, err
	}

	orderCreated, err := be.Orders.CreateOrder(&orders.CreateOrderParams{
		Context: ctx,
		Data: &models.OrderCreateRequest{
			EquipmentID: id,
			RentEnd:     &rentEnd,
			RentStart:   &rentStart,
		},
	}, auth)

	if err != nil {
		return nil, err
	}
	return orderCreated.Payload.ID, nil
}

func checkOrderStatus(ctx context.Context, be *client.Be, auth runtime.ClientAuthInfoWriterFunc, orderId *int64,
	statusName string) (bool, error) {
	orders, err := be.Orders.GetOrdersByStatus(
		orders.NewGetOrdersByStatusParamsWithContext(ctx).WithStatus(statusName), auth)
	if err != nil {
		return false, err
	}
	for _, order := range orders.Payload.Items {
		if *order.ID == *orderId {
			return true, nil
		}
	}
	return false, nil
}

func setParameters(ctx context.Context, client *client.Be, auth runtime.ClientAuthInfoWriterFunc) (*models.Equipment, error) {
	termsOfUse := "https://..."
	cost := int64(3900)
	condition := "удовлетворительное, местами облупляется краска"
	description := "удобная, подойдет для котов любых размеров"
	inventoryNumber := int64(1)

	category, err := client.Categories.GetCategoryByID(categories.NewGetCategoryByIDParamsWithContext(ctx).WithCategoryID(1), auth)
	if err != nil {
		return nil, err
	}

	subCategory, err := client.Subcategories.GetSubcategoryByID(subcategories.NewGetSubcategoryByIDParamsWithContext(ctx).WithSubcategoryID(1), auth)
	if err != nil {
		return nil, err
	}

	location := int64(71)
	mdays := int64(10)
	catName := "Том"
	rDate := int64(1520294400)

	status, err := client.EquipmentStatusName.GetEquipmentStatusName(
		eqStatusName.NewGetEquipmentStatusNameParamsWithContext(ctx).WithStatusID(1), auth)
	if err != nil {
		return nil, err
	}

	f, err := os.Open("../common/cat.jpeg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	petSize, err := client.PetSize.GetAllPetSize(pet_size.NewGetAllPetSizeParamsWithContext(ctx), auth)
	if err != nil {
		return nil, err
	}

	photo, err := client.Photos.CreateNewPhoto(photos.NewCreateNewPhotoParams().WithContext(ctx).WithFile(f), auth)
	if err != nil {
		return nil, err
	}

	cats, err := client.PetKind.GetPetKind(pet_kind.NewGetPetKindParamsWithContext(ctx).WithPetKindID(1), auth)
	if err != nil {
		return nil, err
	}

	supp := "ИП Григорьев Виталий Васильевич"
	techIss := false
	title := "клетка midwest icrate 1"

	var subCategoryInt64 int64
	if subCategory.Payload.Data.ID != nil {
		subCategoryInt64 = *subCategory.Payload.Data.ID
	}

	return &models.Equipment{
		TermsOfUse:       termsOfUse,
		CompensationCost: &cost,
		Condition:        condition,
		Description:      &description,
		InventoryNumber:  &inventoryNumber,
		Category:         category.Payload.Data.ID,
		Subcategory:      subCategoryInt64,
		Location:         &location,
		MaximumDays:      &mdays,
		Name:             &catName,
		NameSubstring:    "box",
		PetKinds:         []int64{*cats.Payload.ID},
		PetSize:          petSize.Payload[0].ID,
		PhotoID:          photo.Payload.Data.ID,
		ReceiptDate:      &rDate,
		Status:           &status.Payload.Data.ID,
		Supplier:         &supp,
		TechnicalIssues:  &techIss,
		Title:            &title,
	}, nil
}

func createEquipment(ctx context.Context, client *client.Be, auth runtime.ClientAuthInfoWriterFunc, model *models.Equipment) (*equipment.CreateNewEquipmentCreated, error) {
	paramsCreate := equipment.NewCreateNewEquipmentParamsWithContext(ctx)
	paramsCreate.NewEquipment = model
	created, err := client.Equipment.CreateNewEquipment(paramsCreate, auth)
	if err != nil {
		return nil, err
	}
	return created, nil
}
