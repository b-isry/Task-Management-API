package Usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "Task-Management/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository represents the repository interface for managing tasks
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetAll(ctx context.Context) ([]*domain.Task, error)
}

// TaskUseCase represents the use case for managing tasks
type TaskUseCase struct {
	repo TaskRepository
}

// NewTaskUseCase creates a new TaskUseCase instance
func NewMockTaskUseCase(repo TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

// CreateTask creates a new task
func (uc *TaskUseCase) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if task.Title == "" {
		return nil, errors.New("task title is required")
	}
	if task.DueDate.Before(time.Now()) {
		return nil, errors.New("due date must be in the future")
	}
	task.Status = domain.StatusPending
	return uc.repo.Create(ctx, task)
}

// GetTaskByID retrieves a task by its ID
func (uc *TaskUseCase) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetTasksByUserID retrieves tasks by user ID
func (uc *TaskUseCase) GetTasksByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Task, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

// GetAllTasks retrieves all tasks
func (uc *TaskUseCase) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	return uc.repo.GetAll(ctx)
}

// UpdateTask updates an existing task
func (uc *TaskUseCase) UpdateTask(ctx context.Context, task *domain.Task) error {
	if task.Title == "" {
		return errors.New("task title is required")
	}
	if task.DueDate.Before(time.Now()) {
		return errors.New("due date must be in the future")
	}
	return uc.repo.Update(ctx, task)
}

// DeleteTask deletes a task by its ID
func (uc *TaskUseCase) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	return uc.repo.Delete(ctx, id)
}

// MockTaskRepository is a mock implementation of the TaskRepository interface.
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Task, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TaskUseCaseTestSuite groups all task use case-related tests
type TaskUseCaseTestSuite struct {
	suite.Suite
	mockRepo *MockTaskRepository
	useCase  *TaskUseCase
}

// SetupSuite runs once before all tests
func (suite *TaskUseCaseTestSuite) SetupSuite() {
	suite.mockRepo = new(MockTaskRepository)
	suite.useCase = NewMockTaskUseCase(suite.mockRepo)
}

// TestCreateTask_Success tests the successful creation of a task
func (suite *TaskUseCaseTestSuite) TestCreateTask_Success() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	suite.mockRepo.On("Create", mock.Anything, task).Return(task, nil)

	result, err := suite.useCase.CreateTask(context.Background(), task)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.StatusPending, result.Status)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateTask_ValidationError tests validation errors during task creation
func (suite *TaskUseCaseTestSuite) TestCreateTask_ValidationError() {
	task := &domain.Task{
		Title:   "",
		DueDate: time.Now().Add(-24 * time.Hour),
	}

	result, err := suite.useCase.CreateTask(context.Background(), task)
	assert.Nil(suite.T(), result)
	assert.EqualError(suite.T(), err, "task title is required")
}

// TestGetTaskByID_Success tests fetching a task by ID successfully
func (suite *TaskUseCaseTestSuite) TestGetTaskByID_Success() {
	taskID := primitive.NewObjectID()
	task := &domain.Task{ID: taskID, Title: "Test Task"}
	suite.mockRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)

	result, err := suite.useCase.GetTaskByID(context.Background(), taskID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), taskID, result.ID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetTasksByUserID_Success tests fetching tasks by user ID successfully
func (suite *TaskUseCaseTestSuite) TestGetTasksByUserID_Success() {
	userID := primitive.NewObjectID()
	tasks := []*domain.Task{
		{ID: primitive.NewObjectID(), Title: "Task 1", UserID: userID},
	}
	suite.mockRepo.On("GetByUserID", mock.Anything, userID).Return(tasks, nil)

	results, err := suite.useCase.GetTasksByUserID(context.Background(), userID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetAllTasks_Success tests fetching all tasks successfully
func (suite *TaskUseCaseTestSuite) TestGetAllTasks_Success() {
	tasks := []*domain.Task{
		{ID: primitive.NewObjectID(), Title: "Task 1"},
		{ID: primitive.NewObjectID(), Title: "Task 2"},
	}
	suite.mockRepo.On("GetAll", mock.Anything).Return(tasks, nil)

	results, err := suite.useCase.GetAllTasks(context.Background())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetTaskByID_NotFound tests fetching a task by ID when it does not exist
func (suite *TaskUseCaseTestSuite) TestGetTaskByID_NotFound() {
	taskID := primitive.NewObjectID()
	suite.mockRepo.On("GetByID", mock.Anything, taskID).Return((*domain.Task)(nil), errors.New("task not found"))

	result, err := suite.useCase.GetTaskByID(context.Background(), taskID)
	assert.Nil(suite.T(), result)
	assert.EqualError(suite.T(), err, "task not found")
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetTasksByUserID_Empty tests fetching tasks by user ID when no tasks exist
func (suite *TaskUseCaseTestSuite) TestGetTasksByUserID_Empty() {
	userID := primitive.NewObjectID()
	suite.mockRepo.On("GetByUserID", mock.Anything, userID).Return([]*domain.Task{}, nil)

	results, err := suite.useCase.GetTasksByUserID(context.Background(), userID)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), results)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUpdateTask_Success tests updating a task successfully
func TestUpdateTask(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	taskUseCase := NewTaskUseCase(mockTaskRepo)

	taskID := primitive.NewObjectID()
	existingTask := &domain.Task{
		ID:      taskID,
		Title:   "Existing Task",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour), // Ensure due date is in the future
	}
	updatedTask := &domain.Task{
		ID:      taskID,
		Title:   "Updated Task",
		Status:  domain.StatusInProgress,
		DueDate: time.Now().Add(48 * time.Hour), // Ensure due date is in the future
	}

	// Mock GetByID call
	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)

	// Mock Update call
	mockTaskRepo.On("Update", mock.Anything, updatedTask).Return(nil)

	// Call UpdateTask
	err := taskUseCase.UpdateTask(context.Background(), updatedTask)

	// Assertions
	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
}

// TestUpdateTask_ValidationError tests validation errors during task update
func (suite *TaskUseCaseTestSuite) TestUpdateTask_ValidationError() {
	task := &domain.Task{
		Title:   "",
		DueDate: time.Now().Add(-24 * time.Hour),
	}

	err := suite.useCase.UpdateTask(context.Background(), task)
	assert.EqualError(suite.T(), err, "task title is required")
}

// TestDeleteTask_Success tests deleting a task successfully
func (suite *TaskUseCaseTestSuite) TestDeleteTask_Success() {
	taskID := primitive.NewObjectID()
	suite.mockRepo.On("Delete", mock.Anything, taskID).Return(nil)

	err := suite.useCase.DeleteTask(context.Background(), taskID)
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateTask_RepositoryError tests repository error during task creation
func (suite *TaskUseCaseTestSuite) TestCreateTask_RepositoryError() {
	task := &domain.Task{
		Title:   "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
	}

	// Mock repository to return an error
	suite.mockRepo.On("Create", mock.Anything, task).Return((*domain.Task)(nil), errors.New("repository error"))

	result, err := suite.useCase.CreateTask(context.Background(), task)

	assert.Nil(suite.T(), result)
	assert.EqualError(suite.T(), err, "repository error")
}

// TestUpdateTask_RepositoryError tests repository error during task update
func (suite *TaskUseCaseTestSuite) TestUpdateTask_RepositoryError() {
	task := &domain.Task{
		ID:      primitive.NewObjectID(),
		Title:   "Updated Task",
		DueDate: time.Now().Add(24 * time.Hour),
	}

	suite.mockRepo.On("Update", mock.Anything, task).Return(errors.New("repository error"))

	err := suite.useCase.UpdateTask(context.Background(), task)

	assert.EqualError(suite.T(), err, "repository error")
}

// Run the test suite
func TestTaskUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseTestSuite))
}
