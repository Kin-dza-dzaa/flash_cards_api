package wordpostgres

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/adrianbrad/psqldocker"
	"github.com/stretchr/testify/suite"
)

// PostgresBase_Suite suite creates throw away docker conainer for tests.
// Docker API should be avialable on 2375 without tls/ssl.
// Doesn't contain any test logic
type PostgresBase_Suite struct {
	suite.Suite
	pg          *Postgres
	pgContainer *psqldocker.Container
	ctx         context.Context
}

// Creates new throw away postgres:alpine container
func (s *PostgresBase_Suite) SetupSuite() {
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
		s.T().Fatalf("PostgresTest - SetupSuite: %s\n", err)
	}

	s.pgContainer = c

	connPool, err := postgres.New(fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s",
		user, pass, c.Port(), db), 10)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.pg = New(connPool)
}

// Container clean-up
func (s *PostgresBase_Suite) TearDownSuite() {
	s.T().Log("clean-up a psql container")
	if err := s.pgContainer.Close(); err != nil {
		s.FailNow(err.Error())
	}
}
