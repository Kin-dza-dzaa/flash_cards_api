package wordpostgres

import (
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/suite"
)

// Sute for testing AddTrans method, embeds PostgresTestBase suite.
type AddTrans_Suite struct {
	PostgresBase_Suite
	tcs []struct {
		Name      string
		WordTrans entity.WordTrans
		WantErr   bool
	}
}

// Sets test case data.
func (s *AddTrans_Suite) SetupTest() {
	s.tcs = []struct {
		Name      string
		WordTrans entity.WordTrans
		WantErr   bool
	}{
		{
			Name: "Add normal word",
			WordTrans: entity.WordTrans{
				Word: "test_word",
			},
			WantErr: false,
		},
		{
			Name:      "Add empty word",
			WordTrans: entity.WordTrans{},
			WantErr:   true,
		},
	}
}

func (s *AddTrans_Suite) Test_AddTrans() {
	for _, tc := range s.tcs {
		s.Run(tc.Name, func() {
			err := s.pg.AddTranslation(s.ctx, tc.WordTrans)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nil")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
		})
	}
}

func Test_AddTrans_Suite(t *testing.T) {
	suite.Run(t, new(AddTrans_Suite))
}
