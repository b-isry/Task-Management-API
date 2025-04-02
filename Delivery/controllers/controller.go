package controllers

import (
	"net/http"

	domain "Task-Management/Domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
}

type UserControllerImpl struct {
	userUseCase domain.UserUseCase
}

type TaskController interface {
	CreateTask(ctx *gin.Context)
	GetTasksByUserID(ctx *gin.Context)
	GetTaskByID(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	DeleteTask(ctx *gin.Context)
	GetAllTasks(ctx *gin.Context)
}

type TaskControllerImpl struct {
	taskUseCase domain.TaskUseCase
}

func NewUserController(userUseCase domain.UserUseCase) *UserControllerImpl {
	return &UserControllerImpl{
		userUseCase: userUseCase,
	}
}

func NewTaskController(taskUseCase domain.TaskUseCase) *TaskControllerImpl {
	return &TaskControllerImpl{
		taskUseCase: taskUseCase,
	}
}

// User Controllers
func (c *UserControllerImpl) Register(ctx *gin.Context) {
	var req domain.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	createdUser, err := c.userUseCase.Register(ctx.Request.Context(), user)
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusConflict, domain.APIResponse{Message: "user already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, domain.APIResponse{Message: "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, domain.APIResponse{
		Message: "User registered successfully",
		Data:    createdUser,
	})
}

func (c *UserControllerImpl) Login(ctx *gin.Context) {
	var req domain.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	user, token, err := c.userUseCase.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Login successful",
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func (c *UserControllerImpl) GetAllUsers(ctx *gin.Context) {
	users, err := c.userUseCase.GetAllUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// Task Controllers
func (c *TaskControllerImpl) CreateTask(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, domain.APIResponse{Message: "unauthorized"})
		return
	}

	var task domain.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: "Invalid user ID"})
		return
	}
	task.UserID = id

	createdTask, err := c.taskUseCase.CreateTask(ctx.Request.Context(), &task)
	if err != nil {
		// Fix: Return 400 for use case errors
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, domain.APIResponse{
		Message: "Task created successfully",
		Data:    createdTask,
	})
}

func (c *TaskControllerImpl) GetTaskByID(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: "Invalid task ID"})
		return
	}

	task, err := c.taskUseCase.GetTaskByID(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, domain.APIResponse{Message: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, domain.APIResponse{Message: "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

func (c *TaskControllerImpl) GetTasksByUserID(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		ctx.JSON(http.StatusUnauthorized, domain.APIResponse{Message: "unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: "Invalid user ID"})
		return
	}

	tasks, err := c.taskUseCase.GetTasksByUserID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Tasks retrieved successfully",
		Data:    tasks,
	})
}

func (c *TaskControllerImpl) GetAllTasks(ctx *gin.Context) {
	tasks, err := c.taskUseCase.GetAllTasks(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Tasks retrieved successfully",
		Data:    tasks,
	})
}

func (c *TaskControllerImpl) UpdateTask(ctx *gin.Context) {
	var task domain.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: "Invalid task ID"})
		return
	}

	task.ID = id
	if err := c.taskUseCase.UpdateTask(ctx.Request.Context(), &task); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Task updated successfully",
	})
}

func (c *TaskControllerImpl) DeleteTask(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: "Invalid task ID"})
		return
	}

	if err := c.taskUseCase.DeleteTask(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.APIResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.APIResponse{
		Message: "Task deleted successfully",
	})
}
