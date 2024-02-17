package main

import (
	"gin-market/controllers"
	"gin-market/infra"
	"gin-market/middleware"
	"gin-market/repository"
	"gin-market/services"
	"gin-market/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	err := infra.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	db := infra.SetUpDB()

	// logging setting
	util.LoggingSettings("./backend/logs/gin_market.log")

	r := setUpRouter(db)

	log.Fatal(r.Run(":8080"))
}

func setUpRouter(db *gorm.DB) *gin.Engine {
	// initialize Gin router
	r := gin.Default()

	// CORS middleware settings only from React origin request
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// save in memory for test
	// itemMemoryRepository := repository.NewItemMemoryRepository(items)

	itemDBRepository := repository.NewItemDBRepository(db)
	itemService := services.NewItemService(itemDBRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	// Grouping
	itemRouter := r.Group("/items")
	itemRouterWithAuth := r.Group("/items", middleware.AuthMiddleware(authService))

	// item router
	itemRouter.GET("", itemController.FindAll)
	itemRouterWithAuth.GET("/mine", itemController.FindMyAll)
	itemRouterWithAuth.GET("/:id", itemController.FindById)
	itemRouterWithAuth.POST("", itemController.Create)
	itemRouterWithAuth.PUT("/:id", itemController.Update)
	itemRouterWithAuth.DELETE("/:id", itemController.Delete)

	// auth router
	authRouter := r.Group("/auth")
	authRouter.POST("/signup", authController.SignUp)
	authRouter.POST("/login", authController.Login)

	authRouter.POST("/logout", authController.Logout)

	return r
}
