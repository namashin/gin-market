package main

import (
	"gin-market/controllers"
	"gin-market/infra"
	"gin-market/middleware"
	"gin-market/repository"
	"gin-market/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
)

func main() {
	infra.Initialize()
	db := infra.SetUpDB()

	// 各パッケージのロギング設定を呼び出す
	LoggingSettings("./logs/main_log/main.log")
	controllers.LoggingSettings("logs/controllers_log/controllers.log")
	services.LoggingSettings("logs/services_log/services.log")
	repository.LoggingSettings("logs/repository_log/repository.log")

	r := setUpRouter(db)

	// サーバーを起動
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func LoggingSettings(logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
}

func setUpRouter(db *gorm.DB) *gin.Engine {
	// [in memory]
	// itemMemoryRepository := repository.NewItemMemoryRepository(items)

	// [in DB]
	itemDBRepository := repository.NewItemDBRepository(db)
	itemService := services.NewItemService(itemDBRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	// Ginのルーターを初期化
	r := gin.Default()

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

	return r
}
