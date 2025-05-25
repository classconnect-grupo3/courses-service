package router

import (
	"courses-service/src/config"
	"courses-service/src/controller"
	"courses-service/src/database"
	"courses-service/src/repository"
	"courses-service/src/service"
	"io"
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
		slog.Error("Failed to create NewRelic application", "error", err)
	}

	r.Use(nrgin.Middleware(app))
}

func InitializeCoursesRoutes(r *gin.Engine, controller *controller.CourseController) {
	r.GET("/courses", controller.GetCourses)
	r.POST("/courses", controller.CreateCourse)
	r.GET("/courses/:id", controller.GetCourseById)
	r.DELETE("/courses/:id", controller.DeleteCourse)
	r.GET("/courses/teacher/:teacherId", controller.GetCourseByTeacherId)
	r.GET("/courses/student/:studentId", controller.GetCoursesByStudentId)
	r.GET("/courses/user/:userId", controller.GetCoursesByUserId)
	r.GET("/courses/title/:title", controller.GetCourseByTitle)
	r.PUT("/courses/:id", controller.UpdateCourse)
}

func InitializeModulesRoutes(r *gin.Engine, controller *controller.ModuleController) {
	r.POST("/modules", controller.CreateModule)
	r.GET("/modules/course/:courseId", controller.GetModulesByCourseId)
	r.GET("/modules/:id", controller.GetModuleById)
	r.DELETE("/modules/:id", controller.DeleteModule)
	r.PUT("/modules/:id", controller.UpdateModule)
}

func InitializeEnrollmentsRoutes(r *gin.Engine, controller *controller.EnrollmentController) {
	r.POST("/courses/:courseId/enroll", controller.EnrollStudent)
	r.POST("/courses/:courseId/unenroll", controller.UnenrollStudent)
}

func NewRouter(config *config.Config) *gin.Engine {
	setUpLogger()
	r := createRouterFromConfig(config)
	addNewRelicMiddleware(r)

	slog.Debug("Connecting to database")

	dbClient, err := database.NewMongoDBClient(config)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	slog.Debug("Connected to database")

	courseRepo := repository.NewCourseRepository(dbClient, config.DBName)
	courseService := service.NewCourseService(courseRepo)
	courseController := controller.NewCourseController(courseService)

	moduleRepo := repository.NewModuleRepository(dbClient, config.DBName)
	moduleService := service.NewModuleService(moduleRepo)
	moduleController := controller.NewModuleController(moduleService)

	enrollmentRepo := repository.NewEnrollmentRepository(dbClient, config.DBName, courseRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo)
	enrollmentController := controller.NewEnrollmentController(enrollmentService)

	InitializeRoutes(r, courseController, moduleController, enrollmentController)
	return r
}

func InitializeRoutes(r *gin.Engine, courseController *controller.CourseController, moduleController *controller.ModuleController, enrollmentController *controller.EnrollmentController) {
	InitializeCoursesRoutes(r, courseController)
	InitializeModulesRoutes(r, moduleController)
	InitializeEnrollmentsRoutes(r, enrollmentController)
}
