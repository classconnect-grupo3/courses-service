package e2e_test

import (
	"courses-service/src/config"
	"courses-service/src/router"
	"courses-service/src/tests/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

var (
	r       = router.NewRouter(config.NewConfig())
	dbSetup *testutil.DBSetup
)

func init() {
	// Initialize database connection for repository tests
	dbSetup = testutil.SetupTestDB()
}

func TestGetEmptyCourses(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	w := httptest.NewRecorder()

	startTime := time.Now()
	endTime := startTime.Add(time.Second * 10)
	course := `{"title": "Test Course", "description": "Test Description", "teacher_id": "123", "capacity": 10, "start_date": "` + startTime.Format(time.RFC3339) + `", "end_date": "` + endTime.Format(time.RFC3339) + `"}`

	req, _ := http.NewRequest("POST", "/courses", strings.NewReader(course))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
