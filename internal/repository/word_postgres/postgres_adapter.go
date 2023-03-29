// Package wordpostgres implements adapter layer for word_api postgres database.
package wordpostgres

import (
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
)

type Postgres struct {
	*postgres.ConnPool
}

func (p *Postgres) Close() {
	p.Pool.Close()
}

func New(pool *postgres.ConnPool) *Postgres {
	return &Postgres{
		pool,
	}
}
