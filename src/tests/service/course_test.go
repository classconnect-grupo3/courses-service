package tests

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockCourseRepository struct{}

func (m *MockCourseRepository) CreateCourse(c model.Course) (*model.Course, error) {
	return &model.Course{
		ID:          primitive.NewObjectID(),
		Title:       c.Title,
		Description: c.Description,
		TeacherUUID: c.TeacherUUID,
		Capacity:    c.Capacity,
	}, nil
}

func (m *MockCourseRepository) GetCourses() ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) GetCourseById(id string) (*model.Course, error) {
	return nil, nil
}

func (m *MockCourseRepository) DeleteCourse(id string) error {
	return nil
}

func (m *MockCourseRepository) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) GetCourseByTitle(title string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func TestCreateCourseWithInvalidCapacity(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.CreateCourse(schemas.CreateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    0,
	})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if course != nil {
		t.Errorf("Expected nil, got %v", course)
	}
}
