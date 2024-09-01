package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
<<<<<<< HEAD

	"github.com/cw2/backend/controllers"
	"github.com/cw2/backend/models"
	"github.com/cw2/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
=======
	"time"

	"github.com/cw2/backend/backup"
	"github.com/cw2/backend/controllers"
	"github.com/cw2/backend/models"
	"github.com/cw2/backend/mw"
	"github.com/cw2/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
>>>>>>> v0.0.2
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
<<<<<<< HEAD
	DB *gorm.DB
	server *gin.Engine
	C controllers.Controller
	R routes.Routes
=======
	DB         *gorm.DB
	server     *gin.Engine
	C          controllers.Controller
	R          routes.Routes
	tradejsloc = "trades/trades.json"
	logger     *logrus.Logger
>>>>>>> v0.0.2
)

func init() {
	var err error
<<<<<<< HEAD
=======
	gin.SetMode(gin.ReleaseMode)

	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logfile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger.Out = logfile
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logfile.Name(),
		MaxSize:    1,    // 单文件最大容量,单位是MB
		MaxBackups: 3,    // 最大保留过期文件个数
		MaxAge:     1,    // 保留过期文件的最大时间间隔,单位是天
		Compress:   true, // 是否需要压缩滚动日志, 使用的 gzip 压缩
	})
	log.SetOutput(logfile)

>>>>>>> v0.0.2
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error is occurred  on .env file please check")
	}
<<<<<<< HEAD
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pwd := os.Getenv("PASSWORD")

	DB,err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
       host, port, user, dbname, pwd)),&gorm.Config{})
=======
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pwd := os.Getenv("DB_PASSWORD")

	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pwd)), &gorm.Config{})
>>>>>>> v0.0.2

	DB.AutoMigrate(
		models.User{},
		&models.Position{},
		&models.Order{},
	)

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
<<<<<<< HEAD
	fmt.Println("🚀 Connected Successfully to the Database")
=======
	log.Println("🚀 Connected Successfully to the Database")
	if _, err = os.Stat("trades"); os.IsNotExist(err) {
		os.Mkdir("trades", 0755)
		os.OpenFile(tradejsloc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	} else if os.Stat(tradejsloc); os.IsNotExist(err) {
		os.OpenFile(tradejsloc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	}
	if _, err = os.Stat("stocks"); os.IsNotExist(err) {
		os.Mkdir("stocks", 0755)
	} else if os.Stat("stocks/marketdata.json"); os.IsNotExist(err) {
		log.Fatal("Please pull the marketdata.")
	}

>>>>>>> v0.0.2
	C = controllers.NewController(DB)
	R = routes.NewRoutes(C)

	server = gin.Default()
<<<<<<< HEAD
}

func main() {
=======
	server.Use(mw.RequestLoggingMiddleware(logger))
}

func main() {

>>>>>>> v0.0.2
	// corsconf := cors.DefaultConfig()
	// corsconf.AllowOrigins = []string{"http://"}
	// corsconf.AllowCredentials = true
	// server.Use(cors.New(corsconf))

	r := server.Group("/api")
<<<<<<< HEAD
	r.GET("/hc", func(ctx *gin.Context) {
		message := "0.0.1"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message":message})
	})
	R.MainR(r)
	log.Fatal(server.Run(":8000"))
}
=======

	r.GET("/hc", func(ctx *gin.Context) {
		message := "0.0.1"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	R.MainR(r)
	go backup.Backup()
	time.Sleep(time.Second * 5)
	go log.Fatal(server.Run(":8000"))

	//Static.UpdateStockTimestamp()

}
>>>>>>> v0.0.2
