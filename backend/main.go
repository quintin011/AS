package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cw2/backend/backup"
	"github.com/cw2/backend/controllers"
	"github.com/cw2/backend/models"
	"github.com/cw2/backend/mw"
	"github.com/cw2/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	server     *gin.Engine
	C          controllers.Controller
	R          routes.Routes
	tradejsloc = "trade/trades.json"
	logger     *logrus.Logger
	sport      string
)

func init() {
	var err error
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

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error is occurred  on .env file please check")
	}
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pwd := os.Getenv("DB_PASSWORD")
	sport = os.Getenv("SERVICE_PORT")

	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pwd)), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	DB.AutoMigrate(
                models.User{},
                &models.Position{},
                &models.Order{},
        )
	log.Println("🚀 Connected Successfully to the Database")
	if _, err = os.Stat("trade"); os.IsNotExist(err) {
		os.Mkdir("trade", 0755)
		os.OpenFile(tradejsloc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	} else if _, err = os.Stat(tradejsloc); os.IsNotExist(err) {
		os.OpenFile(tradejsloc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	}
	if _, err = os.Stat("stocks"); os.IsNotExist(err) {
		os.Mkdir("stocks", 0755)
	} else if _, err = os.Stat("stocks/marketdata.json"); os.IsNotExist(err) {
		log.Fatal("Please pull the marketdata.")
	}

	C = controllers.NewController(DB)
	R = routes.NewRoutes(C)

	server = gin.Default()
	corsconf := cors.DefaultConfig()
	corsconf.ExposeHeaders = []string{"Content-Length", "Content-Type", "Authorization", "X-Uid"}
	corsconf.AllowOrigins = []string{"https://pakatrade.site"}
	corsconf.AllowCredentials = true
	corsconf.AddAllowHeaders("Authorization", "X-Uid")
	//corsconf.AllowHeaders = []string{"Authorization", "X-Uid","Content-Type","Content-Length"}
	// corsconf.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://13.236.191.187:8080", "http://ec2-13-236-191-187.ap-southeast-2.compute.amazonaws.com:8080"}
	// corsconf.AllowCredentials = true
	//server.Use(cors.New(corsconf))
	server.Use(cors.New(corsconf))
	server.Use(mw.RequestLoggingMiddleware(logger))
}

func main() {

	r := server.Group("/api")

	r.GET("/hc", func(ctx *gin.Context) {
		message := "0.0.1"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	R.MainR(r)
	go backup.Backup()
	time.Sleep(time.Second * 5)
	go C.TradeRun()
	time.Sleep(time.Second * 5)
	go log.Fatal(server.RunTLS(":" + sport,"chain.crt","crt.key"))
}
