package wordrepository

import (
	"fmt"
	"os"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/postgres"
	"github.com/adrianbrad/psqldocker"
)

// Creates new throw away postgres:alpine container.
func setupWordRepoContainer(t *testing.T) *WordRepository {
	t.Helper()
	const (
		container = "flash_cards_db_test"
		db        = "test_db"
		user      = "user"
		pass      = "password"
	)

	up, err := os.ReadFile("../../../migrations/up.sql")
	if err != nil {
		t.Fatalf("setupWordRepo - ioutil.ReadFile: %v", err)
	}

	t.Log("starting up a psql container")
	c, err := psqldocker.NewContainer(
		user,
		pass,
		db,
		psqldocker.WithContainerName(container),
		psqldocker.WithSQL(string(up)),
	)
	if err != nil {
		t.Fatalf("setupWordRepo - psqldocker.NewContainer: %v", err)
	}
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Fatalf("setupWordRepo - Cleanup - c.Close: %v", err)
		}
	})

	connPool, err := postgres.New(fmt.Sprintf("postgresql://%s:%s@0.0.0.0:%s/%s", user, pass, c.Port(), db), 10)
	if err != nil {
		t.Fatalf("setupWordRepo - postgres.New: %v", err)
	}

	return New(connPool)
}
