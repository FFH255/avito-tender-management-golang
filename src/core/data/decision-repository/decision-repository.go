package decision_repository

import (
	"tms/src/core/services/repositories"
	"tms/src/pkg/pg"
)

type DecisionRepository struct {
	client pg.Client
}

func New(client pg.Client) repositories.DecisionRepository {
	return DecisionRepository{
		client: client,
	}
}
