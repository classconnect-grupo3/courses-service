package router

import (
	"courses-service/src/config"
	"courses-service/src/controller"
	"courses-service/src/database"
	"courses-service/src/middleware"
	"courses-service/src/repository"
	"courses-service/src/service"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func setUpLogger() {
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

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
	r.GET("/courses/user/:userId", controller.GetCoursesByUserId)
	r.GET("/courses/title/:title", controller.GetCourseByTitle)
	r.GET("/courses/:id", controller.GetCourseById)
	r.DELETE("/courses/:id", controller.DeleteCourse)
	r.PUT("/courses/:id", controller.UpdateCourse)
	r.POST("/courses/:id/aux-teacher/add", controller.AddAuxTeacherToCourse)
	r.DELETE("/courses/:id/aux-teacher/remove", controller.RemoveAuxTeacherFromCourse)
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

	// Esta ruta no requiere autenticación de estudiante
	r.GET("/assignments/:assignmentId/submissions", controller.GetSubmissionsByAssignment)
}

func InitializeEnrollmentsRoutes(r *gin.Engine, controller *controller.EnrollmentController) {
	r.POST("/courses/:id/enroll", controller.EnrollStudent)
	r.POST("/courses/:id/unenroll", controller.UnenrollStudent)
}

func NewRouter(config *config.Config) *gin.Engine {
	setUpLogger()
	r := createRouterFromConfig(config)
	addNewRelicMiddleware(r)

	slog.Debug("Connecting to database")

	dbClient, err := database.NewMongoDBClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	slog.Debug("Connected to database")

	courseRepo := repository.NewCourseRepository(dbClient, config.DBName)
	enrollmentRepo := repository.NewEnrollmentRepository(dbClient, config.DBName, courseRepo)
	assignmentRepository := repository.NewAssignmentRepository(dbClient, config.DBName)
	submissionRepository := repository.NewMongoSubmissionRepository(dbClient.Database(config.DBName))
	moduleRepository := repository.NewModuleRepository(dbClient, config.DBName)

	courseService := service.NewCourseService(courseRepo, enrollmentRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo)
	assignmentService := service.NewAssignmentService(assignmentRepository, courseService)
	submissionService := service.NewSubmissionService(submissionRepository, assignmentRepository, courseService)
	moduleService := service.NewModuleService(moduleRepository)

	courseController := controller.NewCourseController(courseService)
	enrollmentController := controller.NewEnrollmentController(enrollmentService)
	assignmentsController := controller.NewAssignmentsController(assignmentService)
	submissionController := controller.NewSubmissionController(submissionService) // TODO change this when interface is added
	moduleController := controller.NewModuleController(moduleService)

	InitializeRoutes(r, courseController, assignmentsController, submissionController, enrollmentController, moduleController)
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
) {
	InitializeCoursesRoutes(r, courseController)
	InitializeSubmissionRoutes(r, submissionController)
	InitializeAssignmentsRoutes(r, assignmentsController)
	InitializeEnrollmentsRoutes(r, enrollmentController)
	InitializeModulesRoutes(r, moduleController)
}
