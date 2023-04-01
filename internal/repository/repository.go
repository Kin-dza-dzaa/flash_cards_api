// Package repository combines logic of all sub repos.
// Implements adapters layer.
package repository

import (
	googletransrepository "github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/google_trans_repository"
	wordrepository "github.com/Kin-dza-dzaa/flash_cards_api/internal/repository/word_repository"
	googletransclient "github.com/Kin-dza-dzaa/flash_cards_api/pkg/google_trans_client"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
)

type Repository struct {
	*googletransrepository.GoogleTranslate
	*wordrepository.WordRepository
}

func New(client *googletransclient.TranlateClient, pool *postgres.ConnPool, defaultSrcLang string, defaultTrgtLang string) *Repository {
	googleTrans := googletransrepository.New(client, defaultSrcLang, defaultTrgtLang)
	wordPostgres := wordrepository.New(pool)

	return &Repository{
		googleTrans,
		wordPostgres,
	}
}
