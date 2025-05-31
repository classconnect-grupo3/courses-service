package repository_test

type MockEnrollmentRepository struct{}

func (m *MockEnrollmentRepository) IsEnrolled(studentID, courseID string) (bool, error) {
	return true, nil
}

func (m *MockEnrollmentRepository) EnrollStudent(studentID, courseID string) error {
	return nil
}

func (m *MockEnrollmentRepository) UnenrollStudent(studentID, courseID string) error {
	return nil
}
