package tender_repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
	"tms/src/pkg/logger/handlers/slogdiscard"
	"tms/src/pkg/pg"
)

func TestTenderRepository_Get(t *testing.T) {
	client, err := pg.New(slogdiscard.NewDiscardLogger(), pg.Config{
		URL:         "postgres://admin:admin@localhost:5432/tms",
		AutoMigrate: false,
		Migrations:  "",
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := New(*client)

	t.Run("test case", func(t *testing.T) {
		tender, err := repo.Get(context.Background(), repositories.GetTenderDTO{
			ID:             domain.NewID(),
			OrganizationID: nil,
		})
		assert.NoError(t, err)
		assert.NotNil(t, tender)
	})
}
