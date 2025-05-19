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
	r.GET("/courses/title/:title", controller.GetCourseByTitle)
	r.GET("/courses/:id", controller.GetCourseById)
	r.DELETE("/courses/:id", controller.DeleteCourse)
	r.PUT("/courses/:id", controller.UpdateCourse)
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
	
	// Esta ruta no requiere autenticación de estudiante
	r.GET("/assignments/:assignmentId/submissions", controller.GetSubmissionsByAssignment)
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

	courseRepository := repository.NewCourseRepository(dbClient, config.DBName)
	courseService := service.NewCourseService(courseRepository)
	courseController := controller.NewCourseController(courseService)

	assignmentRepository := repository.NewAssignmentRepository(dbClient, config.DBName)
	assignmentService := service.NewAssignmentService(assignmentRepository, courseService)
	assignmentsController := controller.NewAssignmentsController(assignmentService)

	submissionRepository := repository.NewMongoSubmissionRepository(dbClient.Database(config.DBName))
	submissionService := service.NewSubmissionService(submissionRepository, assignmentRepository)
	submissionController := controller.NewSubmissionController(submissionService)

	InitializeRoutes(r, courseController, assignmentsController, submissionController)
	return r
}

func InitializeRoutes(r *gin.Engine, courseController *controller.CourseController, assignmentsController *controller.AssignmentsController, submissionController *controller.SubmissionController) {
	InitializeCoursesRoutes(r, courseController)
	InitializeSubmissionRoutes(r, submissionController)
	InitializeAssignmentsRoutes(r, assignmentsController)
}
