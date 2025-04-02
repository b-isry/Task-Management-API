package routers

import (
	"Task-Management/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userController controllers.UserController,
	taskController controllers.TaskController,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
) *gin.Engine {
	router := gin.Default()

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(authMiddleware)
	{
		// User routes
		protected.GET("/users", userController.GetAllUsers)

		// Task routes
		protected.POST("/tasks", taskController.CreateTask)
		protected.GET("/tasks", taskController.GetTasksByUserID)
		protected.GET("/tasks/:id", taskController.GetTaskByID)
		protected.PUT("/tasks/:id", taskController.UpdateTask)
		protected.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.GET("/tasks", taskController.GetAllTasks)
	}

	return router
}
