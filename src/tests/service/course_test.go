package tests

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if id == "123e4567-e89b-12d3-a456-426614174000" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Test Course",
			Description: "Test Description",
		}, nil
	}
	return nil, nil
}

func (m *MockCourseRepository) DeleteCourse(id string) error {
	return nil
}

func (m *MockCourseRepository) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	if teacherId == "123e4567-e89b-12d3-a456-426614174000" {
		return []*model.Course{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Test Course",
				Description: "Test Description",
				TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
				Capacity:    10,
			},
		}, nil
	}
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) GetCourseByTitle(title string) ([]*model.Course, error) {
	if title == "Test Course" {
		return []*model.Course{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Test Course",
				Description: "Test Description",
				TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
				Capacity:    10,
			},
		}, nil
	}
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error) {
	return &model.Course{
		ID:          primitive.NewObjectID(),
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	}, nil
}

func TestCreateCourseWithInvalidCapacity(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.CreateCourse(schemas.CreateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    0,
	})
	assert.Error(t, err)
	assert.Nil(t, course)
}

func TestCreateCourseWithValidCapacity(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	_, err := courseService.CreateCourse(schemas.CreateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	})
	assert.NoError(t, err)
}

func TestGetCourseById(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.GetCourseById("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.NotNil(t, course)
}

func TestGetCourseByIdWithNonExistentId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.GetCourseById("123e4567-e89b-12d3-a456-426614174001")
	assert.NoError(t, err)
	assert.Nil(t, course)
}

func TestGetCourseByIdWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.GetCourseById("")
	assert.Error(t, err)
	assert.Nil(t, course)
}

func TestGetCourseByTeacherId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	courses, err := courseService.GetCourseByTeacherId("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(courses))
}

func TestGetCourseByTeacherIdWithNonExistentId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	courses, err := courseService.GetCourseByTeacherId("123e4567-e89b-12d3-a456-426614174001")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestGetCourseByTeacherIdWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	courses, err := courseService.GetCourseByTeacherId("")
	assert.Error(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestGetCourseByTitle(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	courses, err := courseService.GetCourseByTitle("Test Course")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(courses))
}

func TestGetCourseByTitleWithNonExistentTitle(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	courses, err := courseService.GetCourseByTitle("Non Existent Title")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestDeleteCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	err := courseService.DeleteCourse("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
}

func TestDeleteCourseWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	err := courseService.DeleteCourse("")
	assert.Error(t, err)
}

func TestUpdateCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.UpdateCourse("123e4567-e89b-12d3-a456-426614174000", schemas.UpdateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	})
	assert.NoError(t, err)
	assert.NotNil(t, course)
}

func TestUpdateCourseWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{})
	course, err := courseService.UpdateCourse("", schemas.UpdateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	})
	assert.Error(t, err)
	assert.Nil(t, course)
}
