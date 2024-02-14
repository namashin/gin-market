package main

import (
	"gin-market/controllers"
	"gin-market/infra"
	"gin-market/middleware"
	"gin-market/repository"
	"gin-market/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	infra.Initialize()
	db := infra.SetUpDB()

	// 各パッケージのロギング設定を呼び出す
	LoggingSettings("./backend/logs/gin_market.log")

	r := setUpRouter(db)

	// サーバーを起動
	log.Fatal(r.Run(":8080"))
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
	// save in memory for test
	// itemMemoryRepository := repository.NewItemMemoryRepository(items)

	itemDBRepository := repository.NewItemDBRepository(db)
	itemService := services.NewItemService(itemDBRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	// Ginのルーターを初期化
	r := gin.Default()

	//config := cors.DefaultConfig()
	//config.AllowHeaders = []string{"Authorization"} // 許可されたヘッダーを指定
	//r.Use(cors.New(config))

	//// CORSミドルウェアを使用して、すべてのオリジンからのリクエストを許可
	// r.Use(cors.Default())
	// CORSミドルウェアを使用して、すべてのオリジンからのリクエストを許可
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// ここでオリジンのチェックを行うこともできます
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

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
