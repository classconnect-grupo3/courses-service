package router

import (
	"courses-service/src/config"
	"courses-service/src/controller"
	"courses-service/src/database"
	"courses-service/src/repository"
	"courses-service/src/service"
	"log"

	"github.com/gin-gonic/gin"
)

func createRouterFromConfig(config *config.Config) *gin.Engine {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	return r
}

func NewRouter(config *config.Config) *gin.Engine {
	r := createRouterFromConfig(config)

	dbClient, err := database.NewMongoDBClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	controller := controller.NewCoursesController(service.NewCourseService(repository.NewCourseRepository(dbClient, config.DBName))) // TODO: dejar esto mas lindo :)
	initializeRoutes(r, controller)
	return r
}

func initializeRoutes(r *gin.Engine, controller *controller.CoursesController) {
	r.GET("/courses", controller.GetCourses)
}
