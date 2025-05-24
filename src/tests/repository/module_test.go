package repository_test

import (
	"courses-service/src/tests/testutil"
	"testing"
)

type ModuleRepositoryMock struct{}

func init() {
	// Initialize database connection for repository tests
	dbSetup = testutil.SetupTestDB()
}

func TestCreateModule(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})
}
