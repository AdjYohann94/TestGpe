package adapter

import (
	"gpe_project/internal/app/adapter/controller"
	"gpe_project/internal/app/adapter/middleware"
	"gpe_project/internal/app/adapter/repository"
	"gpe_project/internal/app/adapter/service"
	"gpe_project/internal/app/application/usecase"

	"github.com/gin-gonic/gin"
)

// Setup boostrap the application by loading and init dependencies.
func Setup() *gin.Engine {
	e := gin.Default()

	// Init database connection
	db := service.GetPostgresqlDB()

	// Auto migration : only on dev mode
	// model.AutoMigrateModels(db)

	// Load repositories and init usecase singleton
	repositories := repository.InitRepositories(db)
	usecases := usecase.InitUsecase(repositories)

	// Enable CORS allow origin
	e.Use(middleware.CORSMiddleware())

	// Active kin validation based on swagger interface
	e.Use(middleware.Kin(service.NewKinValidator()))

	// Non authenticated controllers
	controller.NewPingController(e)
	controller.NewAuthenticationController(e, usecases)
	controller.NewReferentialController(e, usecases)

	// Authentication middleware
	e.Use(middleware.AccessTokenAuthentication())

	// Authenticated controllers
	controller.NewQuizController(e, usecases)
	controller.NewUserController(e, usecases)
	controller.NewBarometerController(e, usecases)
	controller.NewNotificationController(e, usecases)

	return e
}
