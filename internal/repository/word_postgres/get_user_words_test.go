package wordpostgres

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/suite"
)

// Sute for testing GetUserWords method, embeds PostgresTestBase suite.
type GetUserWords_Suite struct {
	PostgresBase_Suite
	tcs []struct {
		Name    string
		Coll    entity.Collection
		Want    *entity.UserWords
		WantErr bool
	}
}

// Sets test case data.
func (s *GetUserWords_Suite) SetupTest() {
	s.tcs = []struct {
		Name    string
		Coll    entity.Collection
		Want    *entity.UserWords
		WantErr bool
	}{
		{
			Name: "Get from empty coll",
			Coll: entity.Collection{
				Name:   "test_coll",
				Word:   "test_word",
				UserID: "12345",
			},
			Want: &entity.UserWords{
				Words: make(map[entity.CollectionName][]entity.WordData, 0),
			},
			WantErr: false,
		},
	}
}

func (s *GetUserWords_Suite) Test_GetUserWords() {
	for _, tc := range s.tcs {
		s.Run(tc.Name, func() {
			actualRes, err := s.pg.GetUserWords(s.ctx, tc.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Want, actualRes, "User words must be as expected")
		})
	}
}

func Test_GetUserWords_Suite(t *testing.T) {
	suite.Run(t, new(GetUserWords_Suite))
}
