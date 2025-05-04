package router

import (
	"courses-service/src/config"
	"courses-service/src/controller"
	"courses-service/src/database"
	"courses-service/src/repository"
	"courses-service/src/service"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
)

func createRouterFromConfig(config *config.Config) *gin.Engine {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	return r
}

func NewRouter(config *config.Config) *gin.Engine {
	r := createRouterFromConfig(config)

	slog.Debug("Connecting to database")

	dbClient, err := database.NewMongoDBClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	slog.Debug("Connected to database")

	controller := controller.NewCoursesController(service.NewCourseService(repository.NewCourseRepository(dbClient, config.DBName))) // TODO: dejar esto mas lindo :)
	initializeRoutes(r, controller)
	return r
}

func initializeRoutes(r *gin.Engine, controller *controller.CoursesController) {
	r.GET("/courses", controller.GetCourses)
	r.POST("/courses", controller.CreateCourse)
	r.GET("/courses/:id", controller.GetCourseById)
	r.DELETE("/courses/:id", controller.DeleteCourse)
	r.GET("/courses/teacher/:teacherId", controller.GetCourseByTeacherId)
	r.GET("/courses/title/:title", controller.GetCourseByTitle)
	r.PUT("/courses/:id", controller.UpdateCourse)
}
