package repositories

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/enttest"
)

type ActiveAreasSuite struct {
	suite.Suite
	ctx         context.Context
	client      *ent.Client
	repository  ActiveAreaRepository
	activeAreas map[int]string
}

func TestActiveAreaSuite(t *testing.T) {
	s := new(ActiveAreasSuite)
	suite.Run(t, s)
}

func (s *ActiveAreasSuite) SetupTest() {
	t := s.T()
	s.ctx = context.Background()
	client := enttest.Open(t, "sqlite3", "file:activeareas?mode=memory&cache=shared&_fk=1")
	s.client = client

	s.activeAreas = make(map[int]string)
	s.activeAreas[1] = "area 1"
	s.activeAreas[2] = "area 2"
	s.activeAreas[3] = "area 3"
	s.activeAreas[4] = "area 4"
	s.activeAreas[5] = "area 5"
	s.activeAreas[6] = "area 6"
	s.activeAreas[7] = "area 7"
	s.activeAreas[8] = "area 8"
	s.activeAreas[9] = "area 9"

	_, err := s.client.ActiveArea.Delete().Exec(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range s.activeAreas {
		_, errCreate := s.client.ActiveArea.Create().SetName(value).Save(s.ctx)
		if errCreate != nil {
			t.Fatal(errCreate)
		}
	}
}

func (s *ActiveAreasSuite) TearDownSuite() {
	s.client.Close()
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreasEmptyOrderBy() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := ""
	orderColumn := "name"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.Nil(t, activeAreas)
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreasEmptyOrderColumn() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "asc"
	orderColumn := ""
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	assert.Error(t, err)
	assert.Nil(t, activeAreas)
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreasOrderByNameDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "desc"
	orderColumn := "name"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.activeAreas), len(activeAreas))
	for i, value := range activeAreas {
		assert.True(t, mapContainsArea(value.Name, s.activeAreas))
		assert.Equal(t, s.activeAreas[len(s.activeAreas)-i], value.Name)
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreasOrderByIDDesc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "desc"
	orderColumn := "id"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.activeAreas), len(activeAreas))
	prevAreaID := math.MaxInt
	for _, value := range activeAreas {
		assert.True(t, mapContainsArea(value.Name, s.activeAreas))
		assert.Less(t, value.ID, prevAreaID)
		prevAreaID = value.ID
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreasOrderByNameAsc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "asc"
	orderColumn := "name"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.activeAreas), len(activeAreas))
	for i, value := range activeAreas {
		assert.True(t, mapContainsArea(value.Name, s.activeAreas))
		assert.Equal(t, s.activeAreas[i+1], value.Name)
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreasOrderByIDAsc() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	orderBy := "asc"
	orderColumn := "id"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.activeAreas), len(activeAreas))
	prevAreaID := -1
	for _, value := range activeAreas {
		assert.True(t, mapContainsArea(value.Name, s.activeAreas))
		assert.Greater(t, value.ID, prevAreaID)
		prevAreaID = value.ID
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_LimitActiveAreas() {
	t := s.T()
	limit := 3
	offset := 0
	orderBy := "asc"
	orderColumn := "name"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(activeAreas))
	for i, value := range activeAreas {
		assert.True(t, mapContainsArea(value.Name, s.activeAreas))
		assert.Equal(t, s.activeAreas[i+1], value.Name)
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_OffsetActiveAreas() {
	t := s.T()
	limit := 6
	offset := 6
	orderBy := "asc"
	orderColumn := "name"
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset, orderBy, orderColumn)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(activeAreas))
	for i, value := range activeAreas {
		assert.True(t, mapContainsArea(value.Name, s.activeAreas))
		assert.Equal(t, s.activeAreas[i+1+offset], value.Name)
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_TotalActiveAreas() {
	t := s.T()
	repository := NewActiveAreaRepository(s.client)
	totalAreas, err := repository.TotalActiveAreas(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.activeAreas), totalAreas)
}

func mapContainsArea(value string, m map[int]string) bool {
	for _, v := range m {
		if value == v {
			return true
		}
	}
	return false
}
