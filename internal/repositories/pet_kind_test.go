package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/enttest"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/middlewares"
)

const petKindEntityName = "pet_kind"

func getClient(t *testing.T, entity string) *ent.Client {
	return enttest.Open(t, "sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", entity))
}

func TestPetKindRepository_Create(t *testing.T) {
	ctx := context.Background()
	client := getClient(t, petKindEntityName)
	name := "test"
	defer client.Close()

	repo := NewPetKindRepository()
	tx, err := client.Tx(ctx)
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	petKind, err := repo.Create(ctx, models.PetKind{Name: &name})
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())

	selectedPetKind, err := client.PetKind.Get(ctx, 1)
	assert.NoError(t, err)

	assert.Equal(t, petKind.ID, selectedPetKind.ID)
	assert.Equal(t, petKind.Name, selectedPetKind.Name)
}

func TestPetKindRepository_GetAll(t *testing.T) {
	ctx := context.Background()
	client := getClient(t, petKindEntityName)
	name := "test"
	_, err := client.PetKind.Create().SetName(name).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	repo := NewPetKindRepository()
	tx, err := client.Tx(ctx)
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	rows, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())

	assert.Equal(t, 1, rows[0].ID)
	assert.Equal(t, name, rows[0].Name)
}

func TestPetKindRepository_Get(t *testing.T) {
	ctx := context.Background()
	client := getClient(t, petKindEntityName)
	name := "test"
	_, err := client.PetKind.Create().SetName(name).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	repo := NewPetKindRepository()
	tx, err := client.Tx(ctx)
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	row, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())

	assert.Equal(t, 1, row.ID)
	assert.Equal(t, name, row.Name)
}

func TestPetKindRepository_Delete(t *testing.T) {
	ctx := context.Background()
	client := getClient(t, petKindEntityName)
	name := "test"
	_, err := client.PetKind.Create().SetName(name).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	repo := NewPetKindRepository()
	tx, err := client.Tx(ctx)
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())

	selected, err := client.PetKind.Get(ctx, 1)
	assert.ErrorContains(t, err, "ent: pet_kind not found")
	assert.Nil(t, selected)
}

func TestPetKindRepository_Update(t *testing.T) {
	ctx := context.Background()
	client := getClient(t, petKindEntityName)
	name := "test"
	name2 := "test2"
	_, err := client.PetKind.Create().SetName(name).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	repo := NewPetKindRepository()
	tx, err := client.Tx(ctx)
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, middlewares.TxContextKey, tx)
	_, err = repo.Update(ctx, 1, &models.PetKind{Name: &name2})
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())

	selected, err := client.PetKind.Get(ctx, 1)
	assert.NoError(t, err)

	assert.Equal(t, 1, selected.ID)
	assert.Equal(t, name2, selected.Name)
}
