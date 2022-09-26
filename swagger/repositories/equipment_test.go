package repositories

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/enttest"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/equipment"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/utils"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/middlewares"
)

type EquipmentSuite struct {
	suite.Suite
	ctx        context.Context
	client     *ent.Client
	repository EquipmentRepository
	equipments map[int]*ent.Equipment
}

func TestEquipmentSuite(t *testing.T) {
	s := new(EquipmentSuite)
	suite.Run(t, s)
}

func (s *EquipmentSuite) SetupTest() {
	t := s.T()
	s.ctx = context.Background()
	client := enttest.Open(t, "sqlite3", "file:equipments?mode=memory&cache=shared&_fk=1")
	s.client = client
	s.repository = NewEquipmentRepository()

	statusName := "status"
	_, err := s.client.Statuses.Delete().Exec(s.ctx) // clean up
	if err != nil {
		t.Fatal(err)
	}
	status, err := s.client.Statuses.Create().SetName(statusName).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	categoryName := "category"
	_, err = s.client.Category.Delete().Exec(s.ctx) // clean up
	if err != nil {
		t.Fatal(err)
	}
	category, err := s.client.Category.Create().SetName(categoryName).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	subcategoryName := "subcategory"
	_, err = s.client.Subcategory.Delete().Exec(s.ctx) // clean up
	if err != nil {
		t.Fatal(err)
	}
	subcategory, err := s.client.Subcategory.Create().SetName(subcategoryName).SetCategory(category).Save(s.ctx)
	if err != nil {
		t.Fatal(err)
	}

	s.equipments = make(map[int]*ent.Equipment)
	s.equipments[1] = &ent.Equipment{
		Name:  "test 1",
		Title: "equipment 1",
	}
	s.equipments[2] = &ent.Equipment{
		Name:  "equipment 2",
		Title: "equipment 2",
	}
	s.equipments[3] = &ent.Equipment{
		Name:  "test 3",
		Title: "equipment 3",
	}
	s.equipments[4] = &ent.Equipment{
		Name:  "equipment 4",
		Title: "equipment 4",
	}
	s.equipments[5] = &ent.Equipment{
		Name:  "test 5",
		Title: "equipment 5",
	}

	_, err = s.client.Equipment.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	for i, value := range s.equipments {
		eq, errCreate := s.client.Equipment.Create().
			SetName(value.Name).SetTitle(value.Title).SetStatus(status).SetCategory(category).SetSubcategory(subcategory).
			Save(s.ctx)
		if errCreate != nil {
			t.Fatal(errCreate)
		}
		s.equipments[i].ID = eq.ID
	}
}

func (s *EquipmentSuite) TearDownSuite() {
	s.client.Close()
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsEmptyOrderBy() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := ""
	orderColumn := equipment.FieldID
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.NoError(t, tx.Rollback())
	assert.Nil(t, equipments)
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsWrongOrderColumn() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.AscOrder
	orderColumn := ""
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.EquipmentsByFilter(ctx, models.EquipmentFilter{},
		limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.NoError(t, tx.Rollback())
	assert.Nil(t, equipments)
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderColumnNotExists() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.AscOrder
	orderColumn := "price"
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.NoError(t, tx.Rollback())
	assert.Nil(t, equipments)
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByIDDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.DescOrder
	orderColumn := equipment.FieldID
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), len(equipments))
	prevEquipmentID := math.MaxInt
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Less(t, value.ID, prevEquipmentID)
		prevEquipmentID = value.ID
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByNameDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.DescOrder
	orderColumn := equipment.FieldName
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), len(equipments))
	prevEquipmentName := "zzzzzzzzzzzzzzzzzzzzzzzzzzz"
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Less(t, value.Name, prevEquipmentName)
		prevEquipmentName = value.Name
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByTitleDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.DescOrder
	orderColumn := equipment.FieldTitle
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), len(equipments))
	prevEquipmentTitle := "zzzzzzzzzzzzzzzzzzzzzzzzzzz"
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Less(t, value.Title, prevEquipmentTitle)
		prevEquipmentTitle = value.Title
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByIDAsc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.AscOrder
	orderColumn := equipment.FieldID
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), len(equipments))
	prevEquipmentID := 0
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Greater(t, value.ID, prevEquipmentID)
		prevEquipmentID = value.ID
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByNameAsc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.AscOrder
	orderColumn := equipment.FieldName
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), len(equipments))
	prevEquipmentName := ""
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Greater(t, value.Name, prevEquipmentName)
		prevEquipmentName = value.Name
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByTitleAsc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := utils.AscOrder
	orderColumn := equipment.FieldTitle
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), len(equipments))
	prevEquipmentTitle := ""
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Greater(t, value.Title, prevEquipmentTitle)
		prevEquipmentTitle = value.Title
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsLimit() {
	t := s.T()
	limit := 3
	offset := 0
	orderBy := utils.AscOrder
	orderColumn := equipment.FieldTitle
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 3, len(equipments))
	for i, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Equal(t, s.equipments[i+1].Name, value.Name)
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOffset() {
	t := s.T()
	limit := 3
	offset := 3
	orderBy := utils.AscOrder
	orderColumn := equipment.FieldTitle
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.AllEquipments(ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 2, len(equipments))
	for i, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Equal(t, s.equipments[i+1+offset].Name, value.Name)
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsTotal() {
	t := s.T()
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	totalEquipment, err := s.repository.AllEquipmentsTotal(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, len(s.equipments), totalEquipment)
}

func (s *EquipmentSuite) TestEquipmentRepository_FindEquipmentsOrderByTitleDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "desc"
	orderColumn := "title"
	filter := models.EquipmentFilter{NameSubstring: "test"}
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.EquipmentsByFilter(ctx, filter, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 3, len(equipments))
	prevEquipmentTitle := "zzzzzzzzzzzzzzzzzzzzzz"
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Contains(t, value.Name, filter.NameSubstring)
		assert.Less(t, value.Title, prevEquipmentTitle)
		prevEquipmentTitle = value.Title
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_FindEquipmentsLimit() {
	t := s.T()
	limit := 2
	offset := 0
	orderBy := "asc"
	orderColumn := "title"
	filter := models.EquipmentFilter{NameSubstring: "test"}
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.EquipmentsByFilter(ctx, filter, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 2, len(equipments))
	prevEquipmentTitle := ""
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Contains(t, value.Name, filter.NameSubstring)
		assert.Greater(t, value.Title, prevEquipmentTitle)
		prevEquipmentTitle = value.Title
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_FindEquipmentsOffset() {
	t := s.T()
	limit := 2
	offset := 2
	orderBy := "asc"
	orderColumn := "name"
	filter := models.EquipmentFilter{NameSubstring: "test"}
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	equipments, err := s.repository.EquipmentsByFilter(ctx, filter, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 1, len(equipments))
	for _, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Contains(t, value.Name, filter.NameSubstring)
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_FindEquipmentsTotal() {
	t := s.T()
	filter := models.EquipmentFilter{NameSubstring: "test"}
	ctx := s.ctx
	tx, err := s.client.Tx(ctx)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	totalEquipment, err := s.repository.EquipmentsByFilterTotal(ctx, filter)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, tx.Commit())
	assert.Equal(t, 3, totalEquipment)
}

func mapContainsEquipment(eq *ent.Equipment, m map[int]*ent.Equipment) bool {
	for _, v := range m {
		if eq.Name == v.Name && eq.Title == v.Title {
			return true
		}
	}
	return false
}
