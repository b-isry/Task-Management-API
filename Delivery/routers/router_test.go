package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserController is a mock implementation of the UserController
type MockUserController struct {
	mock.Mock
}

func (m *MockUserController) Register(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (m *MockUserController) Login(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (m *MockUserController) GetAllUsers(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Users retrieved successfully"})
}

// MockTaskController is a mock implementation of the TaskController
type MockTaskController struct {
	mock.Mock
}

func (m *MockTaskController) CreateTask(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

func (m *MockTaskController) GetTasksByUserID(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Tasks retrieved successfully"})
}

func (m *MockTaskController) GetTaskByID(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Task retrieved successfully"})
}

func (m *MockTaskController) UpdateTask(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (m *MockTaskController) DeleteTask(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (m *MockTaskController) GetAllTasks(ctx *gin.Context) {
	m.Called(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "All tasks retrieved successfully"})
}

// Mock middlewares
func MockAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("user_id", "mockUserID") // Mock user ID
		ctx.Next()
	}
}

func MockAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // Allow all requests
	}
}

// RouterTestSuite groups all router-related tests
type RouterTestSuite struct {
	suite.Suite
	mockUserController *MockUserController
	mockTaskController *MockTaskController
	router             *gin.Engine
}

// SetupSuite runs once before all tests
func (suite *RouterTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest runs before each test
func (suite *RouterTestSuite) SetupTest() {
	suite.mockUserController = new(MockUserController)
	suite.mockTaskController = new(MockTaskController)
	suite.router = SetupRouter(suite.mockUserController, suite.mockTaskController, MockAuthMiddleware(), MockAdminMiddleware())
}

// Test Register Route
func (suite *RouterTestSuite) TestRegisterRoute() {
	suite.mockUserController.On("Register", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodPost, "/api/register", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	suite.mockUserController.AssertExpectations(suite.T())
}

// Test Login Route
func (suite *RouterTestSuite) TestLoginRoute() {
	suite.mockUserController.On("Login", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodPost, "/api/login", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockUserController.AssertExpectations(suite.T())
}

// Test Create Task Route
func (suite *RouterTestSuite) TestCreateTaskRoute() {
	suite.mockTaskController.On("CreateTask", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodPost, "/api/tasks", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	suite.mockTaskController.AssertExpectations(suite.T())
}

// Test Get Tasks By User ID Route
func (suite *RouterTestSuite) TestGetTasksByUserIDRoute() {
	suite.mockTaskController.On("GetTasksByUserID", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskController.AssertExpectations(suite.T())
}

// Test Get Task By ID Route
func (suite *RouterTestSuite) TestGetTaskByIDRoute() {
	suite.mockTaskController.On("GetTaskByID", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/tasks/123", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskController.AssertExpectations(suite.T())
}

// Test Update Task Route
func (suite *RouterTestSuite) TestUpdateTaskRoute() {
	suite.mockTaskController.On("UpdateTask", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodPut, "/api/tasks/123", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskController.AssertExpectations(suite.T())
}

// Test Delete Task Route
func (suite *RouterTestSuite) TestDeleteTaskRoute() {
	suite.mockTaskController.On("DeleteTask", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodDelete, "/api/tasks/123", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskController.AssertExpectations(suite.T())
}

// Test Get All Tasks Route
func (suite *RouterTestSuite) TestGetAllTasksRoute() {
	suite.mockTaskController.On("GetAllTasks", mock.Anything).Return().Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/admin/tasks", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockTaskController.AssertExpectations(suite.T())
}

// Run the test suite
func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
