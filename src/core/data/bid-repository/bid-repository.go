package bid_repository

import (
	"tms/src/core/services/repositories"
	"tms/src/pkg/pg"
)

type BidRepository struct {
	client pg.Client
}

func New(client pg.Client) repositories.BidRepository {
	return BidRepository{
		client: client,
	}
}
