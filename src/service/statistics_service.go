package service

import (
	"bytes"
	"context"
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"encoding/csv"
	"log"
	"strconv"
	"time"
)

// StatisticsService implements the StatisticsServiceInterface
type StatisticsService struct {
	courseRepo     repository.CourseRepositoryInterface
	assignmentRepo repository.AssignmentRepositoryInterface
	enrollmentRepo repository.EnrollmentRepositoryInterface
	submissionRepo repository.SubmissionRepositoryInterface
	forumRepo      repository.ForumRepositoryInterface
}

// NewStatisticsService creates a new instance of StatisticsService
func NewStatisticsService(
	courseRepo repository.CourseRepositoryInterface,
	assignmentRepo repository.AssignmentRepositoryInterface,
	enrollmentRepo repository.EnrollmentRepositoryInterface,
	submissionRepo repository.SubmissionRepositoryInterface,
	forumRepo repository.ForumRepositoryInterface,
) StatisticsServiceInterface {
	return &StatisticsService{
		courseRepo:     courseRepo,
		assignmentRepo: assignmentRepo,
		enrollmentRepo: enrollmentRepo,
		submissionRepo: submissionRepo,
		forumRepo:      forumRepo,
	}
}

// ExportCourseStatsCSV generates a CSV for course statistics and returns the CSV bytes and filename
func (s *StatisticsService) ExportCourseStatsCSV(
	ctx context.Context,
	courseID string,
	from, to time.Time,
) ([]byte, string, error) {
	stats, err := s.GetCourseStatistics(ctx, courseID, from, to)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Comma = ';'

	header := []string{
		"course_id", "course_name", "period_from", "period_to",
		"average_score", "assignment_completion_rate", "exam_completion_rate", "homework_completion_rate",
		"exam_average", "homework_average",
		"total_students", "total_assignments", "total_amount_of_exams", "total_amount_of_homeworks",
		"forum_participation_rate", "forum_unique_participants",
	}

	record := []string{
		stats.CourseID,
		stats.CourseName,
		stats.Period.From.Format("2006-01-02"),
		stats.Period.To.Format("2006-01-02"),
		fmtFloat(stats.AverageScore),
		fmtFloat(stats.AssignmentCompletion),
		fmtFloat(stats.ExamCompletionRate),
		fmtFloat(stats.HomeworkCompletionRate),
		fmtFloat(stats.ExamScoreAverage),
		fmtFloat(stats.HomeworkScoreAverage),
		strconv.Itoa(stats.TotalStudents),
		strconv.Itoa(stats.TotalAssignments),
		strconv.Itoa(stats.TotalAmountOfExams),
		strconv.Itoa(stats.TotalAmountOfHomeworks),
		fmtFloat(stats.ForumParticipationRate),
		strconv.Itoa(stats.ForumUniqueParticipants),
	}

	writer.Write(header)
	writer.Write(record)
	writer.Flush()

	filename := "course_stats_" + courseID + ".csv"
	return buf.Bytes(), filename, nil
}

// ExportStudentStatsCSV generates a CSV for a single student's statistics and returns the CSV bytes and filename
func (s *StatisticsService) ExportStudentStatsCSV(
	ctx context.Context,
	studentID string,
	courseID string,
	from, to time.Time,
) ([]byte, string, error) {
	// Get course details for date defaults if courseID is provided
	var filteredAssignments []*model.Assignment
	var courseTitle string

	course, err := s.courseRepo.GetCourseById(courseID)
	if err != nil {
		return nil, "", err
	}
	if from.IsZero() {
		from = course.StartDate
	}
	if to.IsZero() {
		to = course.EndDate
	}
	assignments, err := s.assignmentRepo.GetAssignmentsByCourseId(courseID)
	if err != nil {
		return nil, "", err
	}
	filteredAssignments, _, _ = s.filterAndSeparateAssignments(assignments, from, to)
	courseTitle = course.Title

	// Get student statistics
	studentStats := s.GetStudentStatistics(ctx, studentID, courseID, filteredAssignments)

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Comma = ';'

	header := []string{
		"student_id", "course_id", "course_name", "period_from", "period_to",
		"average_score", "completion_rate", "participation_rate",
		"completed_assignments", "exam_score", "exam_completed", "homework_score", "homework_completed",
		"forum_posts", "forum_participated", "forum_questions", "forum_answers",
	}

	record := []string{
		studentID,
		courseID,
		courseTitle,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
		fmtFloat(studentStats.PerformanceSummary.AverageScore),
		fmtFloat(studentStats.PerformanceSummary.CompletionRate),
		fmtFloat(studentStats.PerformanceSummary.ParticipationRate),
		strconv.Itoa(studentStats.CompletedAssignments),
		fmtFloat(studentStats.ExamScore),
		strconv.Itoa(studentStats.ExamCompleted),
		fmtFloat(studentStats.HomeworkScore),
		strconv.Itoa(studentStats.HomeworkCompleted),
		strconv.Itoa(studentStats.ForumPosts),
		strconv.FormatBool(studentStats.ForumParticipated),
		strconv.Itoa(studentStats.ForumQuestions),
		strconv.Itoa(studentStats.ForumAnswers),
	}

	writer.Write(header)
	writer.Write(record)
	writer.Flush()

	filename := "student_stats_" + studentID + ".csv"
	return buf.Bytes(), filename, nil
}

// ExportTeacherCoursesStatsCSV generates a CSV with stats for all courses of a teacher
func (s *StatisticsService) ExportTeacherCoursesStatsCSV(
	ctx context.Context,
	teacherID string,
	from, to time.Time,
) ([]byte, string, error) {
	// Get all courses for the teacher
	courses, err := s.courseRepo.GetCourseByTeacherId(teacherID)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Comma = ';'

	header := []string{
		"course_id", "course_name", "period_from", "period_to",
		"average_score", "assignment_completion_rate", "exam_completion_rate", "homework_completion_rate",
		"exam_average", "homework_average",
		"total_students", "total_assignments", "total_amount_of_exams", "total_amount_of_homeworks",
		"forum_participation_rate", "forum_unique_participants",
	}
	writer.Write(header)

	for _, course := range courses {
		stats, err := s.GetCourseStatistics(ctx, course.ID.Hex(), from, to)
		if err != nil {
			// Optionally skip this course or write an error row; here we skip
			continue
		}
		record := []string{
			stats.CourseID,
			stats.CourseName,
			stats.Period.From.Format("2006-01-02"),
			stats.Period.To.Format("2006-01-02"),
			fmtFloat(stats.AverageScore),
			fmtFloat(stats.AssignmentCompletion),
			fmtFloat(stats.ExamCompletionRate),
			fmtFloat(stats.HomeworkCompletionRate),
			fmtFloat(stats.ExamScoreAverage),
			fmtFloat(stats.HomeworkScoreAverage),
			strconv.Itoa(stats.TotalStudents),
			strconv.Itoa(stats.TotalAssignments),
			strconv.Itoa(stats.TotalAmountOfExams),
			strconv.Itoa(stats.TotalAmountOfHomeworks),
			fmtFloat(stats.ForumParticipationRate),
			strconv.Itoa(stats.ForumUniqueParticipants),
		}
		writer.Write(record)
	}
	writer.Flush()

	filename := "teacher_courses_stats_" + teacherID + ".csv"
	return buf.Bytes(), filename, nil
}

//// ------------------------------------------------***********------------------------------------------------ ////
//// ----------------------------------------------HELPER FUNCTIONS------------------------------------------------ ////
//// ------------------------------------------------***********------------------------------------------------ ////

// Retrieves statistics for a specific course
func (s *StatisticsService) GetCourseStatistics(
	ctx context.Context,
	courseID string,
	from, to time.Time,
) (*schemas.CourseStatisticsResponse, error) {
	// Get course details first to access start and end dates
	course, err := s.courseRepo.GetCourseById(courseID)
	if err != nil {
		return nil, err
	}

	// Set default time range if not provided - use course start and end dates
	if from.IsZero() {
		from = course.StartDate
	}
	if to.IsZero() {
		to = course.EndDate
	}

	// Get all enrollments for the course
	enrollments, err := s.enrollmentRepo.GetEnrollmentsByCourseId(courseID)
	if err != nil {
		return nil, err
	}

	// Get all assignments for the course
	assignments, err := s.assignmentRepo.GetAssignmentsByCourseId(courseID)
	if err != nil {
		return nil, err
	}

	// Filter assignments by date and separate by type
	filteredAssignments, examAssignments, homeworkAssignments := s.filterAndSeparateAssignments(assignments, from, to)

	// Process each enrollment to gather student data
	allStudentStats := make([]schemas.StudentStats, 0, len(enrollments))
	for _, enrollment := range enrollments {
		studentStats := s.GetStudentStatistics(ctx, enrollment.StudentID, courseID, filteredAssignments)
		allStudentStats = append(allStudentStats, studentStats)
	}

	// Aggregate stats
	totalScore := 0.0
	totalItems := 0
	examScore := 0.0
	examCompleted := 0
	homeworkScore := 0.0
	homeworkCompleted := 0
	completedCount := 0
	forumParticipantsCount := 0

	for _, studentStats := range allStudentStats {
		totalScore += studentStats.PerformanceSummary.AverageScore
		totalItems++
		examScore += studentStats.ExamScore
		examCompleted += studentStats.ExamCompleted
		homeworkScore += studentStats.HomeworkScore
		homeworkCompleted += studentStats.HomeworkCompleted
		completedCount += studentStats.CompletedAssignments
		if studentStats.ForumParticipated {
			forumParticipantsCount++
		}
	}

	// Calculate course average
	courseAvg := 0.0
	if totalItems > 0 {
		courseAvg = totalScore / float64(totalItems)
	}

	// Calculate exam average
	examAvg := 0.0
	if examCompleted > 0 {
		examAvg = examScore / float64(examCompleted)
	}

	// Calculate homework average
	homeworkAvg := 0.0
	if homeworkCompleted > 0 {
		homeworkAvg = homeworkScore / float64(homeworkCompleted)
	}

	// Calculate exam completion rate
	examCompletionRate := 0.0
	if len(enrollments)*len(examAssignments) > 0 {
		examCompletionRate = float64(examCompleted) / float64(len(enrollments)*len(examAssignments)) * 100
	}

	// Calculate homework completion rate
	homeworkCompletionRate := 0.0
	if len(enrollments)*len(homeworkAssignments) > 0 {
		homeworkCompletionRate = float64(homeworkCompleted) / float64(len(enrollments)*len(homeworkAssignments)) * 100
	}

	// Calculate assignment completion rate across all students
	assignmentCompletionRate := 0.0
	if len(enrollments)*len(filteredAssignments) > 0 {
		assignmentCompletionRate = float64(completedCount) / float64(len(enrollments)*len(filteredAssignments)) * 100
	}

	// Calculate forum participation rate
	forumParticipationRate := 0.0
	if len(enrollments) > 0 {
		forumParticipationRate = float64(forumParticipantsCount) / float64(len(enrollments)) * 100
	}

	return &schemas.CourseStatisticsResponse{
		CourseID:   courseID,
		CourseName: course.Title,
		Period: schemas.Period{
			From: from,
			To:   to,
		},
		AverageScore:            courseAvg,
		AssignmentCompletion:    assignmentCompletionRate,
		ExamCompletionRate:      examCompletionRate,
		HomeworkCompletionRate:  homeworkCompletionRate,
		ExamScoreAverage:        examAvg,
		HomeworkScoreAverage:    homeworkAvg,
		TotalStudents:           len(enrollments),
		TotalAssignments:        len(filteredAssignments),
		TotalAmountOfExams:      len(examAssignments),
		TotalAmountOfHomeworks:  len(homeworkAssignments),
		ForumParticipationRate:  forumParticipationRate,
		ForumUniqueParticipants: forumParticipantsCount,
	}, nil
}

// Helper function to format float with 2 decimals
func fmtFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// calculateStatsForStudent calculates all statistics for a single student
func (s *StatisticsService) GetStudentStatistics(
	ctx context.Context,
	studentID string,
	courseID string,
	filteredAssignments []*model.Assignment,
) schemas.StudentStats {
	studentScore := 0.0
	completedAssignments := 0
	examScore := 0.0
	examCompleted := 0
	homeworkScore := 0.0
	homeworkCompleted := 0

	// Get submissions for this student
	submissions, err := s.submissionRepo.GetByStudent(ctx, studentID)
	if err != nil {
		log.Printf("Error getting submissions for student %s: %v", studentID, err)
		return schemas.StudentStats{}
	}

	// Calculate student statistics
	for _, submission := range submissions {
		// Only count submissions for this course's assignments
		for _, assignment := range filteredAssignments {
			if submission.AssignmentID == assignment.ID.Hex() {
				// Count completion (submitted but not necessarily graded)
				if submission.Status != model.SubmissionStatusDraft {
					completedAssignments++

					// Count type-specific completion
					switch assignment.Type {
					case "exam":
						examCompleted++
					case "homework":
						homeworkCompleted++
					}
				}

				// Calculate scores (only for graded submissions)
				if submission.Score != nil {
					// Safe handling of nil Score
					studentScore += *submission.Score

					// Add to type-specific score totals
					switch assignment.Type {
					case "exam":
						examScore += *submission.Score
					case "homework":
						homeworkScore += *submission.Score
					}
				}
			}
		}
	}

	// Calculate student average
	studentAvg := 0.0
	if completedAssignments > 0 {
		studentAvg = studentScore / float64(completedAssignments)
	}

	// Calculate completion rate
	completionRate := 0.0
	if len(filteredAssignments) > 0 {
		completionRate = float64(completedAssignments) / float64(len(filteredAssignments))
	}

	// Get forum participation for this student
	forumQuestions := 0
	forumAnswers := 0
	questions, _ := s.forumRepo.GetQuestionsByCourseId(courseID)
	for _, question := range questions {
		if question.AuthorID == studentID {
			forumQuestions++
		}
		for _, answer := range question.Answers {
			if answer.AuthorID == studentID {
				forumAnswers++
			}
		}
	}
	forumParticipated := forumQuestions > 0 || forumAnswers > 0

	// Calculate participation rate (simplified)
	participationRate := 0.0
	if len(filteredAssignments) > 0 {
		participationFactors := float64(completedAssignments + forumQuestions + forumAnswers)
		totalFactors := float64(len(filteredAssignments) * 3) // Simple weighting
		participationRate = participationFactors / totalFactors
		if participationRate > 1.0 {
			participationRate = 1.0 // Cap at 100%
		}
	}

	performanceSummary := schemas.StudentPerformanceSummary{
		StudentID:         studentID,
		AverageScore:      studentAvg,
		CompletionRate:    completionRate * 100,    // Convert to percentage
		ParticipationRate: participationRate * 100, // Convert to percentage
	}

	forumPosts := forumQuestions + forumAnswers

	return schemas.StudentStats{
		PerformanceSummary:   performanceSummary,
		StudentScore:         studentScore,
		CompletedAssignments: completedAssignments,
		ExamScore:            examScore,
		ExamCompleted:        examCompleted,
		HomeworkScore:        homeworkScore,
		HomeworkCompleted:    homeworkCompleted,
		ForumPosts:           forumPosts,
		ForumParticipated:    forumParticipated,
		ForumQuestions:       forumQuestions,
		ForumAnswers:         forumAnswers,
	}
}

// filterAndSeparateAssignments filters assignments by date range and separates them by type
func (s *StatisticsService) filterAndSeparateAssignments(assignments []*model.Assignment, from, to time.Time) ([]*model.Assignment, []*model.Assignment, []*model.Assignment) {
	filteredAssignments := []*model.Assignment{}
	examAssignments := []*model.Assignment{}
	homeworkAssignments := []*model.Assignment{}

	for _, assignment := range assignments {
		// Check if assignment is within the date range
		if (assignment.DueDate.After(from) || assignment.DueDate.Equal(from)) &&
			(assignment.DueDate.Before(to) || assignment.DueDate.Equal(to)) {
			filteredAssignments = append(filteredAssignments, assignment)

			// Separate by type
			switch assignment.Type {
			case "exam":
				examAssignments = append(examAssignments, assignment)
			case "homework":
				homeworkAssignments = append(homeworkAssignments, assignment)
			}
		}
	}
	return filteredAssignments, examAssignments, homeworkAssignments
}

// GetBackofficeStatistics returns general system statistics for backoffice
func (s *StatisticsService) GetBackofficeStatistics(ctx context.Context) (*schemas.BackofficeStatisticsResponse, error) {
	// Get general counts
	totalCourses, err := s.courseRepo.CountCourses()
	if err != nil {
		return nil, err
	}

	totalAssignments, err := s.assignmentRepo.CountAssignments()
	if err != nil {
		return nil, err
	}

	totalSubmissions, err := s.submissionRepo.CountSubmissions(ctx)
	if err != nil {
		return nil, err
	}

	totalEnrollments, err := s.enrollmentRepo.CountEnrollments()
	if err != nil {
		return nil, err
	}

	totalForumQuestions, err := s.forumRepo.CountQuestions()
	if err != nil {
		return nil, err
	}

	totalForumAnswers, err := s.forumRepo.CountAnswers()
	if err != nil {
		return nil, err
	}

	// Get course statistics by status
	activeCourses, err := s.courseRepo.CountActiveCourses()
	if err != nil {
		return nil, err
	}

	finishedCourses, err := s.courseRepo.CountFinishedCourses()
	if err != nil {
		return nil, err
	}

	// Get assignment statistics by type
	totalExams, err := s.assignmentRepo.CountAssignmentsByType("exam")
	if err != nil {
		return nil, err
	}

	totalHomeworks, err := s.assignmentRepo.CountAssignmentsByType("homework")
	if err != nil {
		return nil, err
	}

	totalQuizzes, err := s.assignmentRepo.CountAssignmentsByType("quiz")
	if err != nil {
		return nil, err
	}

	// Get submission statistics by status
	draftSubmissions, err := s.submissionRepo.CountSubmissionsByStatus(ctx, model.SubmissionStatusDraft)
	if err != nil {
		return nil, err
	}

	submittedSubmissions, err := s.submissionRepo.CountSubmissionsByStatus(ctx, model.SubmissionStatusSubmitted)
	if err != nil {
		return nil, err
	}

	lateSubmissions, err := s.submissionRepo.CountSubmissionsByStatus(ctx, model.SubmissionStatusLate)
	if err != nil {
		return nil, err
	}

	// Get enrollment statistics by status
	activeEnrollments, err := s.enrollmentRepo.CountEnrollmentsByStatus(model.EnrollmentStatusActive)
	if err != nil {
		return nil, err
	}

	droppedEnrollments, err := s.enrollmentRepo.CountEnrollmentsByStatus(model.EnrollmentStatusDropped)
	if err != nil {
		return nil, err
	}

	completedEnrollments, err := s.enrollmentRepo.CountEnrollmentsByStatus(model.EnrollmentStatusCompleted)
	if err != nil {
		return nil, err
	}

	// Get forum statistics by status
	openForumQuestions, err := s.forumRepo.CountQuestionsByStatus(model.QuestionStatusOpen)
	if err != nil {
		return nil, err
	}

	resolvedForumQuestions, err := s.forumRepo.CountQuestionsByStatus(model.QuestionStatusResolved)
	if err != nil {
		return nil, err
	}

	closedForumQuestions, err := s.forumRepo.CountQuestionsByStatus(model.QuestionStatusClosed)
	if err != nil {
		return nil, err
	}

	// Get teacher and student statistics
	totalUniqueTeachers, err := s.courseRepo.CountUniqueTeachers()
	if err != nil {
		return nil, err
	}

	totalUniqueAuxTeachers, err := s.courseRepo.CountUniqueAuxTeachers()
	if err != nil {
		return nil, err
	}

	totalUniqueStudents, err := s.enrollmentRepo.CountUniqueStudents()
	if err != nil {
		return nil, err
	}

	// Calculate averages
	averageStudentsPerCourse := float64(0)
	if totalCourses > 0 {
		averageStudentsPerCourse = float64(totalEnrollments) / float64(totalCourses)
	}

	averageAssignmentsPerCourse := float64(0)
	if totalCourses > 0 {
		averageAssignmentsPerCourse = float64(totalAssignments) / float64(totalCourses)
	}

	averageSubmissionsPerAssignment := float64(0)
	if totalAssignments > 0 {
		averageSubmissionsPerAssignment = float64(totalSubmissions) / float64(totalAssignments)
	}

	// Get monthly statistics
	coursesCreatedThisMonth, err := s.courseRepo.CountCoursesCreatedThisMonth()
	if err != nil {
		return nil, err
	}

	assignmentsCreatedThisMonth, err := s.assignmentRepo.CountAssignmentsCreatedThisMonth()
	if err != nil {
		return nil, err
	}

	submissionsThisMonth, err := s.submissionRepo.CountSubmissionsThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	enrollmentsThisMonth, err := s.enrollmentRepo.CountEnrollmentsThisMonth()
	if err != nil {
		return nil, err
	}

	return &schemas.BackofficeStatisticsResponse{
		TotalCourses:                    int(totalCourses),
		TotalAssignments:                int(totalAssignments),
		TotalSubmissions:                int(totalSubmissions),
		TotalEnrollments:                int(totalEnrollments),
		TotalForumQuestions:             int(totalForumQuestions),
		TotalForumAnswers:               int(totalForumAnswers),
		ActiveCourses:                   int(activeCourses),
		FinishedCourses:                 int(finishedCourses),
		TotalExams:                      int(totalExams),
		TotalHomeworks:                  int(totalHomeworks),
		TotalQuizzes:                    int(totalQuizzes),
		DraftSubmissions:                int(draftSubmissions),
		SubmittedSubmissions:            int(submittedSubmissions),
		LateSubmissions:                 int(lateSubmissions),
		ActiveEnrollments:               int(activeEnrollments),
		DroppedEnrollments:              int(droppedEnrollments),
		CompletedEnrollments:            int(completedEnrollments),
		OpenForumQuestions:              int(openForumQuestions),
		ResolvedForumQuestions:          int(resolvedForumQuestions),
		ClosedForumQuestions:            int(closedForumQuestions),
		TotalUniqueTeachers:             int(totalUniqueTeachers),
		TotalUniqueAuxTeachers:          int(totalUniqueAuxTeachers),
		TotalUniqueStudents:             int(totalUniqueStudents),
		AverageStudentsPerCourse:        averageStudentsPerCourse,
		AverageAssignmentsPerCourse:     averageAssignmentsPerCourse,
		AverageSubmissionsPerAssignment: averageSubmissionsPerAssignment,
		CoursesCreatedThisMonth:         int(coursesCreatedThisMonth),
		AssignmentsCreatedThisMonth:     int(assignmentsCreatedThisMonth),
		SubmissionsThisMonth:            int(submissionsThisMonth),
		EnrollmentsThisMonth:            int(enrollmentsThisMonth),
	}, nil
}

// GetBackofficeCoursesStats returns detailed course statistics for backoffice
func (s *StatisticsService) GetBackofficeCoursesStats(ctx context.Context) (*schemas.BackofficeCoursesStatsResponse, error) {
	totalCourses, err := s.courseRepo.CountCourses()
	if err != nil {
		return nil, err
	}

	// Get top teachers by course count
	topTeachers, err := s.courseRepo.GetTopTeachersByCourseCount(10)
	if err != nil {
		return nil, err
	}

	// Get courses by status
	activeCourses, err := s.courseRepo.CountActiveCourses()
	if err != nil {
		return nil, err
	}

	finishedCourses, err := s.courseRepo.CountFinishedCourses()
	if err != nil {
		return nil, err
	}

	coursesByStatus := map[string]int{
		"active":   int(activeCourses),
		"finished": int(finishedCourses),
	}

	// Get recent courses
	recentCourses, err := s.courseRepo.GetRecentCourses(10)
	if err != nil {
		return nil, err
	}

	return &schemas.BackofficeCoursesStatsResponse{
		TotalCourses:    int(totalCourses),
		TopTeachers:     topTeachers,
		CoursesByStatus: coursesByStatus,
		RecentCourses:   recentCourses,
	}, nil
}

// GetBackofficeAssignmentsStats returns detailed assignment statistics for backoffice
func (s *StatisticsService) GetBackofficeAssignmentsStats(ctx context.Context) (*schemas.BackofficeAssignmentsStatsResponse, error) {
	totalAssignments, err := s.assignmentRepo.CountAssignments()
	if err != nil {
		return nil, err
	}

	// Get assignments by type
	exams, err := s.assignmentRepo.CountAssignmentsByType("exam")
	if err != nil {
		return nil, err
	}

	homeworks, err := s.assignmentRepo.CountAssignmentsByType("homework")
	if err != nil {
		return nil, err
	}

	quizzes, err := s.assignmentRepo.CountAssignmentsByType("quiz")
	if err != nil {
		return nil, err
	}

	assignmentsByType := map[string]int{
		"exam":     int(exams),
		"homework": int(homeworks),
		"quiz":     int(quizzes),
	}

	// Get assignments by status
	published, err := s.assignmentRepo.CountAssignmentsByStatus("published")
	if err != nil {
		return nil, err
	}

	draft, err := s.assignmentRepo.CountAssignmentsByStatus("draft")
	if err != nil {
		return nil, err
	}

	assignmentsByStatus := map[string]int{
		"published": int(published),
		"draft":     int(draft),
	}

	// Get assignment distribution
	assignmentDistribution, err := s.assignmentRepo.GetAssignmentDistribution()
	if err != nil {
		return nil, err
	}

	// Get recent assignments
	recentAssignments, err := s.assignmentRepo.GetRecentAssignments(10)
	if err != nil {
		return nil, err
	}

	return &schemas.BackofficeAssignmentsStatsResponse{
		TotalAssignments:       int(totalAssignments),
		AssignmentsByType:      assignmentsByType,
		AssignmentsByStatus:    assignmentsByStatus,
		AssignmentDistribution: assignmentDistribution,
		RecentAssignments:      recentAssignments,
	}, nil
}
