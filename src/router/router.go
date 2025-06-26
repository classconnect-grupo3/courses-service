package router

import (
	"courses-service/src/ai"
	"courses-service/src/config"
	"courses-service/src/controller"
	"courses-service/src/database"
	"courses-service/src/middleware"
	"courses-service/src/queues"
	"courses-service/src/repository"
	"courses-service/src/service"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func createRouterFromConfig(config *config.Config) *gin.Engine {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	return r
}

func addNewRelicMiddleware(r *gin.Engine) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("courses-service"),
		newrelic.ConfigLicense("35988c9ba24331e549191b23c94a4cb2FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		log.Fatalf("Failed to create NewRelic application: %v", err)
	}

	r.Use(nrgin.Middleware(app))
}

func InitializeCoursesRoutes(r *gin.Engine, controller *controller.CourseController) {
	r.GET("/courses", controller.GetCourses)
	r.POST("/courses", controller.CreateCourse)
	r.GET("/courses/teacher/:teacherId", controller.GetCourseByTeacherId)
	r.GET("/courses/student/:studentId", controller.GetCoursesByStudentId)
	r.GET("/courses/student/:studentId/favourite", controller.GetFavouriteCourses)
	r.GET("/courses/user/:userId", controller.GetCoursesByUserId)
	r.GET("/courses/title/:title", controller.GetCourseByTitle)
	r.GET("/courses/:id", controller.GetCourseById)
	r.GET("/courses/:id/members", controller.GetCourseMembers)
	r.DELETE("/courses/:id", controller.DeleteCourse)
	r.PUT("/courses/:id", controller.UpdateCourse)
	r.POST("/courses/:id/aux-teacher/add", controller.AddAuxTeacherToCourse)
	r.DELETE("/courses/:id/aux-teacher/remove", controller.RemoveAuxTeacherFromCourse)
	r.POST("/courses/:id/feedback", controller.CreateCourseFeedback)
	r.PUT("/courses/:id/feedback", controller.GetCourseFeedback) // has to be a put because get doesnt receive a body and it was made to receive a body
	r.GET("/courses/:id/feedback/summary", controller.GetCourseFeedbackSummary)
}

func InitializeTeacherActivityRoutes(r *gin.Engine, controller *controller.TeacherActivityController) {
	r.GET("/activity-logs/course/:courseId", controller.GetCourseActivityLogs)
}

func InitializeModulesRoutes(r *gin.Engine, controller *controller.ModuleController) {
	r.POST("/modules", controller.CreateModule)
	r.GET("/modules/course/:courseId", controller.GetModulesByCourseId)
	r.GET("/modules/:id", controller.GetModuleById)
	r.DELETE("/modules/:id", controller.DeleteModule)
	r.PUT("/modules/:id", controller.UpdateModule)
}

func InitializeAssignmentsRoutes(r *gin.Engine, controller *controller.AssignmentsController) {
	r.GET("/assignments", controller.GetAssignments)
	r.POST("/assignments", controller.CreateAssignment)
	r.GET("/assignments/course/:courseId", controller.GetAssignmentsByCourseId)
	r.GET("/assignments/:assignmentId", controller.GetAssignmentById)
	r.PUT("/assignments/:assignmentId", controller.UpdateAssignment)
	r.DELETE("/assignments/:assignmentId", controller.DeleteAssignment)
}

func InitializeSubmissionRoutes(r *gin.Engine, controller *controller.SubmissionController) {
	// Aplicar el middleware de autenticación de estudiantes
	studentAuthGroup := r.Group("")
	studentAuthGroup.Use(middleware.StudentAuth())

	studentAuthGroup.POST("/assignments/:assignmentId/submissions", controller.CreateSubmission)
	studentAuthGroup.GET("/assignments/:assignmentId/submissions/:id", controller.GetSubmission)
	studentAuthGroup.PUT("/assignments/:assignmentId/submissions/:id", controller.UpdateSubmission)
	studentAuthGroup.POST("/assignments/:assignmentId/submissions/:id/submit", controller.SubmitSubmission)
	studentAuthGroup.GET("/students/:studentUUID/submissions", controller.GetSubmissionsByStudent)

	// Aplicar el middleware de autenticación de docentes para calificar
	teacherAuthGroup := r.Group("")
	teacherAuthGroup.Use(middleware.TeacherAuth())
	teacherAuthGroup.PUT("/assignments/:assignmentId/submissions/:id/grade", controller.GradeSubmission)
	teacherAuthGroup.GET("/assignments/:assignmentId/submissions/:id/feedback-summary", controller.GenerateFeedbackSummary)

	// Esta ruta no requiere autenticación de estudiante
	r.GET("/assignments/:assignmentId/submissions", controller.GetSubmissionsByAssignment)
}

func InitializeEnrollmentsRoutes(r *gin.Engine, controller *controller.EnrollmentController) {
	r.GET("/courses/:id/enrollments", controller.GetEnrollmentsByCourseId)
	r.POST("/courses/:id/enroll", controller.EnrollStudent)
	r.DELETE("/courses/:id/unenroll", controller.UnenrollStudent)
	r.POST("/courses/:id/favourite", controller.SetFavouriteCourse)
	r.DELETE("/courses/:id/favourite", controller.UnsetFavouriteCourse)
	r.POST("/courses/:id/student-feedback", controller.CreateFeedback)
	r.PUT("/feedback/student/:id", controller.GetFeedbackByStudentId)
	r.GET("/feedback/student/:id/summary", controller.GetStudentFeedbackSummary)

	// Aplicar el middleware de autenticación de docentes para aprobar estudiantes
	teacherAuthGroup := r.Group("")
	teacherAuthGroup.Use(middleware.TeacherAuth())
	teacherAuthGroup.PUT("/courses/:id/students/:studentId/approve", controller.ApproveStudent)
	teacherAuthGroup.PUT("/courses/:id/students/:studentId/disapprove", controller.DisapproveStudent)
}

func InitializeForumRoutes(r *gin.Engine, controller *controller.ForumController) {
	// Question endpoints
	r.POST("/forum/questions", controller.CreateQuestion)
	r.GET("/forum/questions/:questionId", controller.GetQuestionById)
	r.GET("/forum/courses/:courseId/questions", controller.GetQuestionsByCourseId)
	r.PUT("/forum/questions/:questionId", controller.UpdateQuestion)
	r.DELETE("/forum/questions/:questionId", controller.DeleteQuestion)

	// Answer endpoints
	r.POST("/forum/questions/:questionId/answers", controller.AddAnswer)
	r.PUT("/forum/questions/:questionId/answers/:answerId", controller.UpdateAnswer)
	r.DELETE("/forum/questions/:questionId/answers/:answerId", controller.DeleteAnswer)
	r.POST("/forum/questions/:questionId/answers/:answerId/accept", controller.AcceptAnswer)

	// Vote endpoints
	r.POST("/forum/questions/:questionId/vote", controller.VoteQuestion)
	r.POST("/forum/questions/:questionId/answers/:answerId/vote", controller.VoteAnswer)
	r.DELETE("/forum/questions/:questionId/vote", controller.RemoveVoteFromQuestion)
	r.DELETE("/forum/questions/:questionId/answers/:answerId/vote", controller.RemoveVoteFromAnswer)

	// Search endpoints
	r.GET("/forum/courses/:courseId/search", controller.SearchQuestions)

	// Forum participants endpoints
	r.GET("/forum/courses/:courseId/participants", controller.GetForumParticipants)
}

// InitializeStatisticsRoutes sets up all statistics-related routes
func InitializeStatisticsRoutes(r *gin.Engine, controller *controller.StatisticsController) {
	// Group all statistics routes
	statisticsGroup := r.Group("/statistics")

	// Apply teacher authentication middleware for all statistics routes
	// Only teachers should be able to access these analytics
	statisticsGroup.Use(middleware.TeacherAuth())

	// Course statistics endpoints - supports JSON or CSV via ?format=csv
	statisticsGroup.GET("/courses/:courseId", controller.GetCourseStatistics)

	// Student statistics endpoints - supports JSON or CSV via ?format=csv
	statisticsGroup.GET("/students/:studentId", controller.GetStudentStatistics)

	// Teacher's courses statistics endpoint
	statisticsGroup.GET("/teachers/:teacherId/courses", controller.GetTeacherCoursesStatistics)

	// Backoffice statistics endpoints
	backofficeGroup := r.Group("/backoffice/statistics")

	// General system statistics
	backofficeGroup.GET("/general", controller.GetBackofficeStatistics)

	// Detailed course statistics
	backofficeGroup.GET("/courses", controller.GetBackofficeCoursesStats)

	// Detailed assignment statistics
	backofficeGroup.GET("/assignments", controller.GetBackofficeAssignmentsStats)
}

func NewRouter(config *config.Config) *gin.Engine {
	r := createRouterFromConfig(config)
	addNewRelicMiddleware(r)

	slog.Debug("Connecting to database")

	dbClient, err := database.NewMongoDBClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	slog.Debug("Connected to database")

	aiClient := ai.NewAiClient(config)
	notificationsQueue, err := queues.NewNotificationsQueue(config)
	if err != nil {
		log.Fatalf("Failed to create notifications queue: %v", err)
	}

	courseRepo := repository.NewCourseRepository(dbClient, config.DBName)
	enrollmentRepo := repository.NewEnrollmentRepository(dbClient, config.DBName, courseRepo)
	assignmentRepository := repository.NewAssignmentRepository(dbClient, config.DBName)
	submissionRepository := repository.NewMongoSubmissionRepository(dbClient.Database(config.DBName))
	moduleRepository := repository.NewModuleRepository(dbClient, config.DBName)
	forumRepository := repository.NewForumRepository(dbClient, config.DBName)
	activityLogRepo := repository.NewTeacherActivityLogRepository(dbClient, config.DBName)

	courseService := service.NewCourseService(courseRepo, enrollmentRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo, submissionRepository)
	assignmentService := service.NewAssignmentService(assignmentRepository, courseService)
	submissionService := service.NewSubmissionService(submissionRepository, assignmentRepository, courseService, aiClient)
	moduleService := service.NewModuleService(moduleRepository)
	forumService := service.NewForumService(forumRepository, courseRepo)
	statisticsService := service.NewStatisticsService(courseRepo, assignmentRepository, enrollmentRepo, submissionRepository, forumRepository)
	activityService := service.NewTeacherActivityService(activityLogRepo, courseRepo)

	courseController := controller.NewCourseController(courseService, aiClient, activityService)
	enrollmentController := controller.NewEnrollmentController(enrollmentService, aiClient, activityService)
	assignmentsController := controller.NewAssignmentsController(assignmentService, notificationsQueue, activityService)
	submissionController := controller.NewSubmissionController(submissionService, notificationsQueue, activityService, assignmentService)
	moduleController := controller.NewModuleController(moduleService, activityService)
	forumController := controller.NewForumController(forumService, activityService)
	statisticsController := controller.NewStatisticsController(statisticsService)
	activityController := controller.NewTeacherActivityController(activityService, courseService)

	InitializeRoutes(r, courseController, assignmentsController, submissionController, enrollmentController, moduleController, forumController, statisticsController, activityController)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // endpoint to consult the swagger documentation
	return r
}

func InitializeRoutes(
	r *gin.Engine,
	courseController *controller.CourseController,
	assignmentsController *controller.AssignmentsController,
	submissionController *controller.SubmissionController,
	enrollmentController *controller.EnrollmentController,
	moduleController *controller.ModuleController,
	forumController *controller.ForumController,
	statisticsController *controller.StatisticsController,
	activityController *controller.TeacherActivityController,
) {
	InitializeCoursesRoutes(r, courseController)
	InitializeSubmissionRoutes(r, submissionController)
	InitializeAssignmentsRoutes(r, assignmentsController)
	InitializeEnrollmentsRoutes(r, enrollmentController)
	InitializeModulesRoutes(r, moduleController)
	InitializeForumRoutes(r, forumController)
	InitializeStatisticsRoutes(r, statisticsController)
	InitializeTeacherActivityRoutes(r, activityController)
}
