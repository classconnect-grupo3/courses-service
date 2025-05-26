package e2e_test

import (
	"courses-service/src/config"
	"courses-service/src/router"
	"courses-service/src/schemas"
	"courses-service/src/tests/testutil"
	"encoding/json"
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

func TestGetCoursesForAUserThatIsTeacherAndStudent(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("enrollments")
	})

	teacherId1 := "123"
	teacherId2 := "456"
	studentId := "789"

	// Create first course
	w := httptest.NewRecorder()
	course1JSON := `{"title": "Test Course 1", "description": "Test Description 1", "teacher_id": "` + teacherId1 + `", "capacity": 10, "start_date": "` + time.Now().Format(time.RFC3339) + `", "end_date": "` + time.Now().Add(time.Second*10).Format(time.RFC3339) + `"}`
	req, _ := http.NewRequest("POST", "/courses", strings.NewReader(course1JSON))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Extract course1 ID
	var course1Response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &course1Response)
	assert.Equal(t, nil, err)
	course1ID := course1Response["id"].(string)

	// Create second course
	w = httptest.NewRecorder()
	course2JSON := `{"title": "Test Course 2", "description": "Test Description 2", "teacher_id": "` + teacherId2 + `", "capacity": 10, "start_date": "` + time.Now().Format(time.RFC3339) + `", "end_date": "` + time.Now().Add(time.Second*10).Format(time.RFC3339) + `"}`
	req, _ = http.NewRequest("POST", "/courses", strings.NewReader(course2JSON))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Extract course2 ID
	var course2Response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &course2Response)
	assert.Equal(t, nil, err)
	course2ID := course2Response["id"].(string)

	// Create third course (where student is teacher)
	w = httptest.NewRecorder()
	course3JSON := `{"title": "Test Course 3", "description": "Test Description 3", "teacher_id": "` + studentId + `", "capacity": 10, "start_date": "` + time.Now().Format(time.RFC3339) + `", "end_date": "` + time.Now().Add(time.Second*10).Format(time.RFC3339) + `"}`
	req, _ = http.NewRequest("POST", "/courses", strings.NewReader(course3JSON))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Extract course3 ID
	var course3Response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &course3Response)
	assert.Equal(t, nil, err)
	// We don't need to store course3ID since we're not enrolling the student in it
	// but verifying it as a teacher course

	// Enroll student in course1
	w = httptest.NewRecorder()
	enrollment1JSON := `{"student_id": "` + studentId + `"}`
	req, _ = http.NewRequest("POST", "/courses/"+course1ID+"/enroll", strings.NewReader(enrollment1JSON))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Enroll student in course2
	w = httptest.NewRecorder()
	enrollment2JSON := `{"student_id": "` + studentId + `"}`
	req, _ = http.NewRequest("POST", "/courses/"+course2ID+"/enroll", strings.NewReader(enrollment2JSON))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Get courses for student
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/courses/user/"+studentId, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response schemas.GetCoursesByUserIdResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, nil, err)

	// Verify teacher courses (studentId is also a teacher for course3)
	assert.Equal(t, 1, len(response.Teacher))
	assert.Equal(t, "Test Course 3", response.Teacher[0].Title)
	assert.Equal(t, studentId, response.Teacher[0].TeacherUUID)

	// Verify student courses (enrolled in course1 and course2)
	assert.Equal(t, 2, len(response.Student))

	// Sort the student courses by title for consistent assertions
	studentCourses := make(map[string]string)
	for _, course := range response.Student {
		studentCourses[course.Title] = course.ID.Hex()
	}

	// Check course IDs match what we expect
	assert.Equal(t, true, studentCourses["Test Course 1"] == course1ID || studentCourses["Test Course 2"] == course2ID)
}
