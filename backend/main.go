package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	// "time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	// "github.com/cw2/backend/backup"
	"github.com/cw2/backend/controllers"
	"github.com/cw2/backend/models"
	"github.com/cw2/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/natefinch/lumberjack"
	// "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	server     *gin.Engine
	C          controllers.Controller
	R          routes.Routes
	tradejsloc = "trade/trades.json"
	ginLambda *ginadapter.GinLambda
)

func init() {
	var err error
	gin.SetMode(gin.ReleaseMode)
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pwd := os.Getenv("DB_PASSWORD")

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

	C = controllers.NewController(DB)
	R = routes.NewRoutes(C)

	server = gin.Default()
	corsconf := cors.DefaultConfig()
	corsconf.ExposeHeaders = []string{"Content-Length", "Content-Type", "X-Uid","X-Amzn-Remapped-Authorization"}
	//corsconf.AllowOrigins = []string{"https://pakatrade.site"}
	corsconf.AllowAllOrigins = true
	corsconf.AllowCredentials = true
	corsconf.AddAllowHeaders("Authorization", "X-Uid")

	server.Use(cors.New(corsconf))
	r := server.Group("/api")
	r.GET("/hc", func(ctx *gin.Context) {
		message := "0.0.1"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	R.MainR(r)
	ginLambda = ginadapter.New(server)
}

func main() {
	lambda.Start(GinRequestHandler)
}

func GinRequestHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {	
	return ginLambda.ProxyWithContext(ctx, req)
}
