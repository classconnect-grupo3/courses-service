package tests

import (
	"courses-service/src/tests/testutil"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

var testDB *mongo.Client
var testDBName string
var dbSetup *testutil.DBSetup

func TestMain(m *testing.M) {
	// Set up database connection before running tests
	dbSetup = testutil.SetupTestDB()
	testDB = dbSetup.Client
	testDBName = dbSetup.DBName

	// Run tests
	code := m.Run()

	// Clean up after tests
	dbSetup.CleanupCollection("courses")
	testutil.CleanupTestDB(testDB)

	os.Exit(code)
}
