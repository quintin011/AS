package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cw2/backend/controllers"
	"github.com/cw2/backend/models"
	"github.com/cw2/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	server *gin.Engine
	C controllers.Controller
	R routes.Routes
)

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error is occurred  on .env file please check")
	}
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pwd := os.Getenv("PASSWORD")

	DB,err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
       host, port, user, dbname, pwd)),&gorm.Config{})

	DB.AutoMigrate(
		models.User{},
		&models.Position{},
		&models.Order{},
	)

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
	C = controllers.NewController(DB)
	R = routes.NewRoutes(C)

	server = gin.Default()
}

func main() {
	// corsconf := cors.DefaultConfig()
	// corsconf.AllowOrigins = []string{"http://"}
	// corsconf.AllowCredentials = true
	// server.Use(cors.New(corsconf))

	r := server.Group("/api")
	r.GET("/hc", func(ctx *gin.Context) {
		message := "0.0.1"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message":message})
	})
	R.MainR(r)
	log.Fatal(server.Run(":8000"))
}