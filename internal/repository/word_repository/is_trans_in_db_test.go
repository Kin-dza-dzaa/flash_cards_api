package wordrepository

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/suite"
)

// Sute for testing IsTransInDB_Suite method, embeds PostgresTestBase suite.
type IsTransInDB_Suite struct {
	WordRepository_Base_Suite
	tcs []struct {
		Name    string
		Coll    entity.Collection
		Want    bool
		WantErr bool
	}
}

// Sets test case data and adds trans to DB if needed.
func (s *IsTransInDB_Suite) SetupTest() {
	s.tcs = []struct {
		Name    string
		Coll    entity.Collection
		Want    bool
		WantErr bool
	}{
		{
			Name: "Not existing trans",
			Coll: entity.Collection{
				Name:   "test_coll",
				Word:   "not_exist_word",
				UserID: "12345",
			},
			Want:    false,
			WantErr: false,
		},
		{
			Name: "Existing trans",
			Coll: entity.Collection{
				Name:   "test_coll",
				Word:   "test_word",
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
		}
	}
}

func (s *IsTransInDB_Suite) Test_GetUserWords() {
	for _, tc := range s.tcs {
		s.Run(tc.Name, func() {
			res, err := s.pg.IsTransInDB(s.ctx, tc.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Want, res, "Bool result must be as expected")
		})
	}
}

func Test_IsTransInDB_Suite(t *testing.T) {
	suite.Run(t, new(IsTransInDB_Suite))
}
