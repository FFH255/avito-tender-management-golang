package tender_repository

import (
	"tms/src/core/services/repositories"
	"tms/src/pkg/pg"
)

type TenderRepository struct {
	client pg.Client
}

func New(client pg.Client) repositories.TenderRepository {
	return TenderRepository{
		client: client,
	}
}
