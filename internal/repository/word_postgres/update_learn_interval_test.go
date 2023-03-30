package wordpostgres

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/suite"
)

// Sute for testing UpdateLearnInterval method, embeds PostgresTestBase suite.
type UpdateLearnInterval_Suite struct {
	PostgresBase_Suite
	tcs []struct {
		Name    string
		Coll    entity.Collection
		WantErr bool
	}
}

// Sets test case data.
func (s *UpdateLearnInterval_Suite) SetupTest() {
	s.tcs = []struct {
		Name    string
		Coll    entity.Collection
		WantErr bool
	}{
		{
			Name: "Update word",
			Coll: entity.Collection{
				Name:   "test_coll",
				Word:   "test_word",
				UserID: "12345",
			},
			WantErr: false,
		},
	}
}

func (s *UpdateLearnInterval_Suite) Test_UpdateLearnInterval() {
	for _, tc := range s.tcs {
		s.Run(tc.Name, func() {
			err := s.pg.UpdateLearnInterval(s.ctx, tc.Coll)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func Test_UpdateLearnInterval_Suite(t *testing.T) {
	suite.Run(t, new(UpdateLearnInterval_Suite))
}
