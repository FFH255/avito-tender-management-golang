package bid_repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
	"tms/src/pkg/logger/handlers/slogdiscard"
	"tms/src/pkg/pg"
)

func TestBidRepository_Get(t *testing.T) {
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
		bid, err := repo.Get(context.Background(), repositories.GetBidDTO{
			ID:       domain.NewID(),
			AuthorID: nil,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, bid)
	})
}
