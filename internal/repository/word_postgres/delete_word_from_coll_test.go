package wordpostgres

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/suite"
)

// Sute for testing DeleteWordFromColl method, embeds PostgresTestBase suite.
type DeleteWordFromColl_Suite struct {
	PostgresBase_Suite
	tcs []struct {
		Name    string
		Coll    *entity.Collection
		WantErr bool
	}
}

// Sets test case data.
func (s *DeleteWordFromColl_Suite) SetupTest() {
	s.tcs = []struct {
		Name    string
		Coll    *entity.Collection
		WantErr bool
	}{
		{
			Name: "Delete word",
			Coll: &entity.Collection{
				Name:   "test_coll",
				Word:   "test_word",
				UserID: "12345",
			},
			WantErr: false,
		},
		{
			Name:    "Delete not existing word",
			Coll:    &entity.Collection{},
			WantErr: false,
		},
	}
}

func (s *DeleteWordFromColl_Suite) Test_DeleteWordFromColl() {
	for _, tc := range s.tcs {
		s.Run(tc.Name, func() {
			err := s.pg.DeleteWordFromCollection(s.ctx, tc.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func Test_DeleteWordFromColl_Suite(t *testing.T) {
	suite.Run(t, new(DeleteWordFromColl_Suite))
}
