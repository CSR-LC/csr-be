package repositories

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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

func (s *ActiveAreasSuite) TestActiveAreaRepository_AllActiveAreas() {
	t := s.T()
	limit := math.MaxInt
	offset := 0
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(s.activeAreas), len(activeAreas))
	for _, value := range activeAreas {
		assert.Contains(t, s.activeAreas, value.ID)
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_LimitActiveAreas() {
	t := s.T()
	limit := 3
	offset := 0
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(activeAreas))
	for i, value := range activeAreas {
		assert.Contains(t, s.activeAreas, value.ID)
		assert.Equal(t, s.activeAreas[i+1], value.Name)
	}
}

func (s *ActiveAreasSuite) TestActiveAreaRepository_OffsetActiveAreas() {
	t := s.T()
	limit := 6
	offset := 6
	repository := NewActiveAreaRepository(s.client)
	activeAreas, err := repository.AllActiveAreas(s.ctx, limit, offset)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(activeAreas))
	for i, value := range activeAreas {
		assert.Contains(t, s.activeAreas, value.ID)
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
