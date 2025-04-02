package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Task-Management/Delivery/controllers"
	"Task-Management/Delivery/routers"
	repository "Task-Management/Repository"
	"Task-Management/Usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongoDB() (*mongo.Client, *mongo.Database, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, err
	}

	db := client.Database("taskmanager")
	return client, db, nil
}

func initServer(router http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

func runServer(srv *http.Server, suppressLogs bool) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if !suppressLogs {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()
}

func main() {
	// Initialize MongoDB
	client, db, err := initMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize use cases
	userUseCase := Usecases.NewUserUseCase(userRepo)
	taskUseCase := Usecases.NewTaskUseCase(taskRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userUseCase)
	taskController := controllers.NewTaskController(taskUseCase)

	// Setup router
	// Define middleware functions
	middleware1 := func(c *gin.Context) {
		// Example middleware logic
		log.Println("Middleware1 executed")
		c.Next()
	}

	middleware2 := func(c *gin.Context) {
		// Example middleware logic
		log.Println("Middleware2 executed")
		c.Next()
	}

	// Setup router with middlewares
	router := routers.SetupRouter(userController, taskController, middleware1, middleware2)

	// Initialize and run server
	srv := initServer(router)
	runServer(srv, false)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
