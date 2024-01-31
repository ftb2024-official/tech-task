package main

import (
	"os"

	handler "tech_task/handlers"
	middleware "tech_task/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func connect(user, pswd, dbName string) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: pswd,
		Database: dbName,
	})

	if db == nil {
		logrus.Fatal("Failed to connect to the database...")
	}

	return db
}

func setupRoutes(router *gin.Engine, db *pg.DB) {
	api := router.Group("/api")
	api2 := router.Group("/api2")

	api.Use(middleware.CheckQueryParamMiddleware).GET("/users", handler.GetUsersHandler(db))
	api.Use(middleware.CheckNameMiddleware).POST("/users", handler.CreateUserHandler(db))

	api2.Use(middleware.CheckParamMiddleware).PATCH("/users/:id", handler.EditUserHandler(db))
	api2.Use(middleware.CheckParamMiddleware).DELETE("/users/:id", handler.DeleteUserHandler(db))
}

func main() {
	var log = logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.WarnLevel
	log.Out = os.Stdout

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading '.env' file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPswd := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db := connect(dbUser, dbPswd, dbName)
	defer db.Close()

	router := gin.Default() // initialize Gin router

	setupRoutes(router, db) // define API routes

	if err := router.Run(":8080"); err != nil {
		logrus.Fatal("Unable to run the server at port 8080 :(")
	}
}
