package wordrepository

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/adrianbrad/psqldocker"
	"github.com/stretchr/testify/suite"
)

// WordRepository_Base_Suite suite creates throw away docker conainer for tests.
// Docker API should be avialable on 2375 without tls/ssl.
// Doesn't contain any test logic
type WordRepository_Base_Suite struct {
	suite.Suite
	pg          *WordRepository
	pgContainer *psqldocker.Container
	ctx         context.Context
}

// Creates new throw away postgres:alpine container
func (s *WordRepository_Base_Suite) SetupSuite() {
	const (
		container = "flash_cards_db_test"
		db        = "test_db"
		user      = "user"
		pass      = "password"
	)

	s.ctx = context.Background()

	up, err := ioutil.ReadFile("../../../migrations/up.sql")
	if err != nil {
		s.FailNow(err.Error())
	}

	s.T().Log("starting up a psql container")
	c, err := psqldocker.NewContainer(
		user,
		pass,
		db,
		psqldocker.WithContainerName(container),
		psqldocker.WithSQL(string(up)),
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.pgContainer = c

	connPool, err := postgres.New(fmt.Sprintf("postgresql://%s:%s@0.0.0.0:%s/%s",
		user, pass, c.Port(), db), 10)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.pg = New(connPool)
}

// Container clean-up
func (s *WordRepository_Base_Suite) TearDownSuite() {
	s.T().Log("clean-up a psql container")
	if err := s.pgContainer.Close(); err != nil {
		s.FailNow(err.Error())
	}
}
