package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/enttest"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/middlewares"
)

type orderStatusTestSuite struct {
	suite.Suite
	ctx           context.Context
	client        *ent.Client
	statusNameMap map[int]string
	repository    OrderStatusRepository
	adminUser     *ent.User
	order         *ent.Order
}

func TestOrderStatusSuite(t *testing.T) {
	s := new(orderStatusTestSuite)
	suite.Run(t, s)
}

func (s *orderStatusTestSuite) SetupTest() {
	t := s.T()
	s.ctx = context.Background()
	client := enttest.Open(t, "sqlite3", "file:orderstatus?mode=memory&cache=shared&_fk=1")
	s.client = client

	s.statusNameMap = map[int]string{ // list of statuses. copy of sql migration
		1: "review",
		2: "approved",
		3: "in progress",
		4: "rejected",
		5: "closed",
	}

	_, err := s.client.StatusName.Delete().Exec(s.ctx) // clean up
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.client.User.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.client.Order.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, statusName := range s.statusNameMap { // create statuses
		_, errCreation := s.client.StatusName.Create().SetStatus(statusName).Save(s.ctx)
		if errCreation != nil {
			t.Fatal(errCreation)
		}
	}

	user, err := s.client.User.Create().SetLogin("admin").SetName("user"). // create user
										SetPassword("admin").SetEmail("test@example.com").Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.adminUser = user

	order, err := s.client.Order.Create().SetDescription("test order").SetQuantity(1). // create order
												SetRentStart(time.Now()).SetRentEnd(time.Now()).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.order = order

	s.repository = NewOrderStatusRepository()
}

func (s *orderStatusTestSuite) TearDownSuite() {
	s.client.Close()
}

func (s *orderStatusTestSuite) TestOrderStatusRepository_UpdateStatus() {
	t := s.T()
	userID := s.adminUser.ID
	comment := "test comment"
	createdAt := strfmt.DateTime(time.Now().UTC())
	orderID := int64(s.order.ID)
	status, ok := s.statusNameMap[1]
	if !ok {
		t.Error("cant find status with id 1")
	}
	data := models.NewOrderStatus{
		Comment:   &comment,
		CreatedAt: &createdAt,
		OrderID:   &orderID,
		Status:    &status,
	}
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	err = s.repository.UpdateStatus(ctx, userID, data)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())
	_, err = s.client.OrderStatus.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *orderStatusTestSuite) TestOrderStatusRepository_StatusHistory_Empty() {
	t := s.T()
	orderID := s.order.ID
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	statuses, err := s.repository.StatusHistory(ctx, orderID)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())
	assert.Empty(t, statuses)
}

func (s *orderStatusTestSuite) TestOrderStatusRepository_StatusHistory() {
	t := s.T()
	orderID := s.order.ID
	// create order status

	orderStatus, err := s.client.OrderStatus.Create().SetComment("test comment").SetCurrentDate(time.Now().UTC()).
		SetOrderID(s.order.ID).SetUsersID(s.adminUser.ID).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	statuses, err := s.repository.StatusHistory(ctx, orderID)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 1, len(statuses))
	assert.Equal(t, orderStatus.ID, statuses[0].ID)
	assert.Equal(t, orderStatus.Comment, statuses[0].Comment)
	assert.Equal(t, orderStatus.CurrentDate, statuses[0].CurrentDate)
	_, err = s.client.OrderStatus.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *orderStatusTestSuite) TestOrderStatusRepository_GetOrderCurrentStatus() {
	t := s.T()
	orderID := s.order.ID
	// create order status
	orderStatus, err := s.client.OrderStatus.Create().SetComment("test comment").SetCurrentDate(time.Now().UTC()).
		SetOrderID(s.order.ID).SetUsersID(s.adminUser.ID).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	status, err := s.repository.GetOrderCurrentStatus(ctx, orderID)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())
	assert.Equal(t, orderStatus.ID, status.ID)
	assert.Equal(t, orderStatus.Comment, status.Comment)
	assert.Equal(t, orderStatus.CurrentDate, status.CurrentDate)
	_, err = s.client.OrderStatus.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *orderStatusTestSuite) TestOrderStatusRepository_GetUserStatusHistory() {
	t := s.T()
	userID := s.adminUser.ID
	// create order status
	orderStatus, err := s.client.OrderStatus.Create().SetComment("test comment").SetCurrentDate(time.Now().UTC()).
		SetOrderID(s.order.ID).SetUsersID(s.adminUser.ID).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	statuses, err := s.repository.GetUserStatusHistory(ctx, userID)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 1, len(statuses))
	assert.Equal(t, orderStatus.ID, statuses[0].ID)
	assert.Equal(t, orderStatus.Comment, statuses[0].Comment)
	assert.Equal(t, orderStatus.CurrentDate, statuses[0].CurrentDate)
	_, err = s.client.OrderStatus.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
}
