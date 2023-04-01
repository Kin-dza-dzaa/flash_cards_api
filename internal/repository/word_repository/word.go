// Package wordrepository implements adapter layer for postgres database.
package wordrepository

import (
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
)

type WordRepository struct {
	*postgres.ConnPool
}

func (p *WordRepository) Close() {
	p.Pool.Close()
}

func New(pool *postgres.ConnPool) *WordRepository {
	return &WordRepository{
		pool,
	}
}
