package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/ajianaz/gofin-full/api/tests/integration/testhelpers"
)

// testApp is shared across all integration tests in this package.
// It is initialised once in TestMain and cleaned up after all tests finish.
var testApp *testhelpers.TestApp

func TestMain(m *testing.M) {
	cfg := testhelpers.NewTestConfig()

	app, err := testhelpers.NewTestApp(cfg)
	if err != nil {
		// If the test database is not available, skip integration tests
		// rather than failing. This allows `go test ./...` to pass in
		// environments without Docker/PostgreSQL.
		fmt.Fprintf(os.Stderr, "SKIP integration tests: %v\n", err)
		os.Exit(0)
	}
	testApp = app

	code := m.Run()

	testApp.Cleanup()
	os.Exit(code)
}
