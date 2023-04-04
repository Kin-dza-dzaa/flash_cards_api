// Package service implements application/usecase layer of the app, combines sub services.
package service

import wordservice "github.com/Kin-dza-dzaa/flash_cards_api/internal/service/word_service"

type Service struct {
	*wordservice.WordService
}

func New(dbAdapter wordservice.WordRepository, googletransAdapter wordservice.Translator) *Service {
	ws := wordservice.New(dbAdapter, googletransAdapter)
	return &Service{
		ws,
	}
}
