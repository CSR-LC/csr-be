package repositories

import (
	"context"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/enttest"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

type EquipmentSuite struct {
	suite.Suite
	ctx        context.Context
	client     *ent.Client
	repository EquipmentRepository
	equipments map[int]ent.Equipment
}

func TestEquipmentSuite(t *testing.T) {
	s := new(EquipmentSuite)
	suite.Run(t, s)
}

func (s *EquipmentSuite) SetupTest() {
	t := s.T()
	s.ctx = context.Background()
	client := enttest.Open(t, "sqlite3", "file:activeareas?mode=memory&cache=shared&_fk=1")
	s.client = client

	s.equipments = make(map[int]ent.Equipment)
	s.equipments[1] = ent.Equipment{
		Name:  "equipment 1",
		Title: "equipment 1",
	}
	s.equipments[2] = ent.Equipment{
		Name:  "equipment 2",
		Title: "equipment 2",
	}
	s.equipments[3] = ent.Equipment{
		Name:  "equipment 3",
		Title: "equipment 3",
	}
	s.equipments[4] = ent.Equipment{
		Name:  "equipment 4",
		Title: "equipment 4",
	}
	s.equipments[5] = ent.Equipment{
		Name:  "equipment 5",
		Title: "equipment 5",
	}

	_, err := s.client.Equipment.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range s.equipments {
		_, errCreate := s.client.Equipment.Create().
			SetName(value.Name).SetTitle(value.Title).Save(s.ctx)
		if errCreate != nil {
			t.Fatal(errCreate)
		}
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
	orderColumn := "id"
	repository := NewEquipmentRepository(s.client)
	equipment, err := repository.EquipmentsByFilter(s.ctx, models.EquipmentFilter{},
		limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.Nil(t, equipment)
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsEmptyOrderColumn() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "asc"
	orderColumn := ""
	repository := NewEquipmentRepository(s.client)
	equipment, err := repository.EquipmentsByFilter(s.ctx, models.EquipmentFilter{},
		limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.Nil(t, equipment)
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderColumnNotExists() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "asc"
	orderColumn := "price"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.Nil(t, equipments)
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsOrderByIDDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "desc"
	orderColumn := "id"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "desc"
	orderColumn := "name"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "desc"
	orderColumn := "title"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "asc"
	orderColumn := "id"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "asc"
	orderColumn := "name"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "asc"
	orderColumn := "title"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "asc"
	orderColumn := "name"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
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
	orderBy := "asc"
	orderColumn := "name"
	repository := NewEquipmentRepository(s.client)
	equipments, err := repository.AllEquipments(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(equipments))
	for i, value := range equipments {
		assert.True(t, mapContainsEquipment(value, s.equipments))
		assert.Equal(t, s.equipments[i+1+offset].Name, value.Name)
	}
}

func (s *EquipmentSuite) TestEquipmentRepository_AllEquipmentsTotal() {
	t := s.T()
	repository := NewEquipmentRepository(s.client)
	totalEquipment, err := repository.AllEquipmentsTotal(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.equipments), totalEquipment)
}

func mapContainsEquipment(eq *ent.Equipment, m map[int]ent.Equipment) bool {
	for _, v := range m {
		if eq.Name == v.Name && eq.Title == v.Title {
			return true
		}
	}
	return false
}
