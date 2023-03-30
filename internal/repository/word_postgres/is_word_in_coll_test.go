package wordpostgres

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/suite"
)

// Sute for testing IsWordInColl method, embeds PostgresTestBase suite.
type IsWordInColl_Suite struct {
	PostgresBase_Suite
	tcs []struct {
		Name    string
		Coll    entity.Collection
		Want    bool
		WantErr bool
	}
}

// Sets test case data and adds words to collection if needed.
func (s *IsWordInColl_Suite) SetupTest() {
	s.tcs = []struct {
		Name    string
		Coll    entity.Collection
		Want    bool
		WantErr bool
	}{
		{
			Name:    "Not existing word",
			Coll:    entity.Collection{},
			Want:    false,
			WantErr: false,
		},
		{
			Name: "Existing word",
			Coll: entity.Collection{
				Word:   "some_word",
				Name:   "test_coll",
				UserID: "12345",
			},
			Want:    true,
			WantErr: false,
		},
	}

	for _, tc := range s.tcs {
		if tc.Want {
			if err := s.pg.AddTranslation(s.ctx,
				entity.WordTrans{Word: tc.Coll.Word}); err != nil {
				s.FailNow(err.Error())
			}
			if err := s.pg.AddWordToCollection(s.ctx, tc.Coll); err != nil {
				s.FailNow(err.Error())
			}
		}
	}
}

func (s *IsWordInColl_Suite) Test_GetUserWords() {
	for _, tc := range s.tcs {
		s.Run(tc.Name, func() {
			res, err := s.pg.IsWordInCollection(s.ctx, tc.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Want, res, "Bool result must be as expected")
		})
	}
}

func Test_IsWordInColl_Suite(t *testing.T) {
	suite.Run(t, new(IsWordInColl_Suite))
}
