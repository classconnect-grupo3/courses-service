package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"log/slog"
	"slices"
)

type TeacherActivityService struct {
	activityLogRepo repository.TeacherActivityLogRepositoryInterface
	courseRepo      repository.CourseRepositoryInterface
}

func NewTeacherActivityService(
	activityLogRepo repository.TeacherActivityLogRepositoryInterface,
	courseRepo repository.CourseRepositoryInterface,
) *TeacherActivityService {
	return &TeacherActivityService{
		activityLogRepo: activityLogRepo,
		courseRepo:      courseRepo,
	}
}

// LogActivityIfAuxTeacher logs an activity only if the teacher is an auxiliary teacher for the course
func (s *TeacherActivityService) LogActivityIfAuxTeacher(courseID, teacherUUID, activityType, description string) {
	// Get the course to check if teacher is auxiliary
	course, err := s.courseRepo.GetCourseById(courseID)
	if err != nil {
		slog.Error("Failed to get course for activity logging", "error", err, "courseID", courseID)
		return
	}

	// Check if teacher is auxiliary (not the main teacher)
	if course.TeacherUUID == teacherUUID {
		// This is the main teacher, don't log
		return
	}

	if !slices.Contains(course.AuxTeachers, teacherUUID) {
		// This teacher is not an auxiliary teacher for this course
		return
	}

	// Log the activity
	err = s.activityLogRepo.LogActivity(courseID, teacherUUID, activityType, description)
	if err != nil {
		slog.Error("Failed to log teacher activity", "error", err, "courseID", courseID, "teacherUUID", teacherUUID)
	} else {
		slog.Debug("Auxiliary teacher activity logged", "courseID", courseID, "teacherUUID", teacherUUID, "activityType", activityType)
	}
}

// GetCourseActivityLogs returns all activity logs for a course
func (s *TeacherActivityService) GetCourseActivityLogs(courseID string) ([]*model.TeacherActivityLog, error) {
	return s.activityLogRepo.GetLogsByCourse(courseID)
} 