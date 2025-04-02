package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Task-Management/Domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserUseCase is a mock implementation of the UserUseCase interface
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) GetUserByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserUseCase) Register(ctx context.Context, user *Domain.User) (*Domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserUseCase) Login(ctx context.Context, email, password string) (*Domain.User, string, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*Domain.User), args.String(1), args.Error(2)
}

func (m *MockUserUseCase) GetAllUsers(ctx context.Context) ([]*Domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Domain.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUseCase) UpdateUser(ctx context.Context, user *Domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockTaskUseCase is a mock implementation of the TaskUseCase interface
type MockTaskUseCase struct {
	mock.Mock
}

func (m *MockTaskUseCase) CreateTask(ctx context.Context, task *Domain.Task) (*Domain.Task, error) {
	args := m.Called(ctx, task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskUseCase) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskUseCase) GetTasksByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Domain.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Domain.Task), args.Error(1)
}

func (m *MockTaskUseCase) GetAllTasks(ctx context.Context) ([]*Domain.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Domain.Task), args.Error(1)
}

func (m *MockTaskUseCase) UpdateTask(ctx context.Context, task *Domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskUseCase) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestSuite struct for grouping tests
type ControllerTestSuite struct {
	suite.Suite
	mockUserUseCase *MockUserUseCase
	mockTaskUseCase *MockTaskUseCase
	router          *gin.Engine
}

// SetupSuite runs once before all tests
func (suite *ControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest runs before each test
func (suite *ControllerTestSuite) SetupTest() {
	suite.mockUserUseCase = new(MockUserUseCase)
	suite.mockTaskUseCase = new(MockTaskUseCase)
	suite.router = gin.Default()
}

// Test UserController: Register Success
func (suite *ControllerTestSuite) TestUserController_Register_Success() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.POST("/register", controller.Register)

	mockUser := &Domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     "user",
	}

	suite.mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*Domain.User")).Return(mockUser, nil)

	body, _ := json.Marshal(Domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     "user",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// Test UserController: Login Success
func (suite *ControllerTestSuite) TestUserController_Login_Success() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.POST("/login", controller.Login)

	mockUser := &Domain.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	suite.mockUserUseCase.On("Login", mock.Anything, "john@example.com", "password123").Return(mockUser, "mockToken", nil)

	body, _ := json.Marshal(Domain.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	})

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// Test TaskController: CreateTask Success
func (suite *ControllerTestSuite) TestTaskController_CreateTask_Success() {
	controller := NewTaskController(suite.mockTaskUseCase)

	// Middleware to mock user_id in the context
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID().Hex())
		c.Next()
	})

	suite.router.POST("/tasks", controller.CreateTask)

	mockTask := &Domain.Task{
		Title:       "Test Task",
		Description: "This is a test task",
	}

	suite.mockTaskUseCase.On("CreateTask", mock.Anything, mock.AnythingOfType("*Domain.Task")).Return(mockTask, nil)

	body, _ := json.Marshal(mockTask)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: DeleteTask Success
func (suite *ControllerTestSuite) TestTaskController_DeleteTask_Success() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.DELETE("/tasks/:id", controller.DeleteTask)

	mockID := primitive.NewObjectID()
	suite.mockTaskUseCase.On("DeleteTask", mock.Anything, mockID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+mockID.Hex(), nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test UserController: Register Validation Error
func (suite *ControllerTestSuite) TestUserController_Register_ValidationError() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.POST("/register", controller.Register)

	body := `{"email": "invalid-email", "password": "short", "role": "invalid-role"}`

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test TaskController: CreateTask Use Case Error
func (suite *ControllerTestSuite) TestTaskController_CreateTask_UseCaseError() {
	controller := NewTaskController(suite.mockTaskUseCase)

	// Middleware to mock user_id in the context
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID().Hex()) // Ensure user_id is set
		c.Next()
	})

	suite.router.POST("/tasks", controller.CreateTask)

	suite.mockTaskUseCase.On("CreateTask", mock.Anything, mock.Anything).Return(nil, errors.New("use case error"))

	body := `{"title": "Test Task", "due_date": "2024-12-31T00:00:00Z"}`

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code) // Expect 400
}

// Test UserController: Register with Malformed JSON
func (suite *ControllerTestSuite) TestUserController_Register_MalformedJSON() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.POST("/register", controller.Register)

	body := `{"name": "John Doe", "email": "john@example.com", "password":}` // Malformed JSON

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test UserController: Register Duplicate User
func (suite *ControllerTestSuite) TestUserController_Register_DuplicateUser() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.POST("/register", controller.Register)

	mockError := errors.New("user already exists")
	suite.mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*Domain.User")).Return(nil, mockError)

	body, _ := json.Marshal(Domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     "user",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusConflict, resp.Code) // Fix: Expect 409
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// Test TaskController: CreateTask Unauthorized Access
func (suite *ControllerTestSuite) TestTaskController_CreateTask_Unauthorized() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.POST("/tasks", controller.CreateTask)

	body := `{"title": "Test Task", "description": "This is a test task"}`

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Ensure middleware does not set user_id
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code) // Fix: Expect 401
}

// Test TaskController: GetTask Invalid Task ID
func (suite *ControllerTestSuite) TestTaskController_GetTask_InvalidID() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks/:id", controller.GetTaskByID)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/invalid-id", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test UserController: GetAllUsers Success
func (suite *ControllerTestSuite) TestUserController_GetAllUsers_Success() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.GET("/users", controller.GetAllUsers)

	mockUsers := []*Domain.User{
		{Name: "John Doe", Email: "john@example.com"},
		{Name: "Jane Doe", Email: "jane@example.com"},
	}

	suite.mockUserUseCase.On("GetAllUsers", mock.Anything).Return(mockUsers, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// Test UserController: GetAllUsers Internal Server Error
func (suite *ControllerTestSuite) TestUserController_GetAllUsers_InternalServerError() {
	controller := NewUserController(suite.mockUserUseCase)
	suite.router.GET("/users", controller.GetAllUsers)

	suite.mockUserUseCase.On("GetAllUsers", mock.Anything).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

// Test TaskController: GetAllTasks Success
func (suite *ControllerTestSuite) TestTaskController_GetAllTasks_Success() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks", controller.GetAllTasks)

	mockTasks := []*Domain.Task{
		{Title: "Task 1", Description: "Description 1"},
		{Title: "Task 2", Description: "Description 2"},
	}

	suite.mockTaskUseCase.On("GetAllTasks", mock.Anything).Return(mockTasks, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: GetAllTasks Internal Server Error
func (suite *ControllerTestSuite) TestTaskController_GetAllTasks_InternalServerError() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks", controller.GetAllTasks)

	suite.mockTaskUseCase.On("GetAllTasks", mock.Anything).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

// Test TaskController: UpdateTask Success
func (suite *ControllerTestSuite) TestTaskController_UpdateTask_Success() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.PUT("/tasks/:id", controller.UpdateTask)

	mockID := primitive.NewObjectID()
	mockTask := Domain.Task{Title: "Updated Task", Description: "Updated Description"}
	mockTask.ID = mockID // Ensure the task ID is set

	// Fix: Properly set up the mock to return nil for the UpdateTask call
	suite.mockTaskUseCase.On("UpdateTask", mock.Anything, &mockTask).Return(nil)

	body, _ := json.Marshal(mockTask)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+mockID.Hex(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code) // Expect 200 OK
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: UpdateTask Invalid Task ID
func (suite *ControllerTestSuite) TestTaskController_UpdateTask_InvalidTaskID() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.PUT("/tasks/:id", controller.UpdateTask)

	body := `{"title": "Updated Task", "description": "Updated Description"}`
	req, _ := http.NewRequest(http.MethodPut, "/tasks/invalid-id", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test TaskController: DeleteTask Invalid Task ID
func (suite *ControllerTestSuite) TestTaskController_DeleteTask_InvalidTaskID() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.DELETE("/tasks/:id", controller.DeleteTask)

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/invalid-id", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test TaskController: GetTaskByID Success
func (suite *ControllerTestSuite) TestTaskController_GetTaskByID_Success() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks/:id", controller.GetTaskByID)

	mockID := primitive.NewObjectID()
	mockTask := &Domain.Task{ID: mockID, Title: "Test Task", Description: "Test Description"}

	suite.mockTaskUseCase.On("GetTaskByID", mock.Anything, mockID).Return(mockTask, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+mockID.Hex(), nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: GetTaskByID Not Found
func (suite *ControllerTestSuite) TestTaskController_GetTaskByID_NotFound() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks/:id", controller.GetTaskByID)

	mockID := primitive.NewObjectID()
	suite.mockTaskUseCase.On("GetTaskByID", mock.Anything, mockID).Return(nil, errors.New("task not found"))

	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+mockID.Hex(), nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: GetTasksByUserID Success
func (suite *ControllerTestSuite) TestTaskController_GetTasksByUserID_Success() {
	controller := NewTaskController(suite.mockTaskUseCase)

	// Middleware to mock user_id in the context
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID().Hex())
		c.Next()
	})

	suite.router.GET("/tasks/user", controller.GetTasksByUserID)

	mockTasks := []*Domain.Task{
		{Title: "Task 1", Description: "Description 1"},
		{Title: "Task 2", Description: "Description 2"},
	}

	suite.mockTaskUseCase.On("GetTasksByUserID", mock.Anything, mock.Anything).Return(mockTasks, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/user", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: GetTasksByUserID Invalid UserID
func (suite *ControllerTestSuite) TestTaskController_GetTasksByUserID_InvalidUserID() {
	controller := NewTaskController(suite.mockTaskUseCase)

	// Middleware to mock invalid user_id in the context
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", "invalid-id")
		c.Next()
	})

	suite.router.GET("/tasks/user", controller.GetTasksByUserID)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/user", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test TaskController: Internal Server Error
func (suite *ControllerTestSuite) TestTaskController_InternalServerError() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks/:id", controller.GetTaskByID)

	mockID := primitive.NewObjectID()
	suite.mockTaskUseCase.On("GetTaskByID", mock.Anything, mockID).Return(nil, errors.New("internal server error"))

	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+mockID.Hex(), nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code) // Expect 500
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Test TaskController: Unauthorized Access
func (suite *ControllerTestSuite) TestTaskController_UnauthorizedAccess() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks/user", controller.GetTasksByUserID)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/user", nil)
	resp := httptest.NewRecorder()

	// Ensure middleware does not set user_id
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", nil) // Explicitly set user_id to nil
		c.Next()
	})

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code) // Expect 401
}

// Test TaskController: Bad Request Error
func (suite *ControllerTestSuite) TestTaskController_BadRequestError() {
	controller := NewTaskController(suite.mockTaskUseCase)
	suite.router.GET("/tasks/:id", controller.GetTaskByID)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/invalid-id", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test UserController: Invalid User ID
func (suite *ControllerTestSuite) TestUserController_InvalidUserID() {
	controller := NewTaskController(suite.mockTaskUseCase)

	// Middleware to mock invalid user_id in the context
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", "invalid-id")
		c.Next()
	})

	suite.router.GET("/tasks/user", controller.GetTasksByUserID)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/user", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

// Test TaskController: Bad Request on Task Creation
func (suite *ControllerTestSuite) TestTaskController_CreateTask_BadRequest() {
	controller := NewTaskController(suite.mockTaskUseCase)

	// Middleware to mock user_id in the context
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID().Hex()) // Ensure user_id is set
		c.Next()
	})

	suite.router.POST("/tasks", controller.CreateTask)

	body := `{"title": "Test Task", "description":}` // Malformed JSON

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code) // Expect 400
}

// Run the test suite
func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
