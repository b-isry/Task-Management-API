package repository

import (
	"context"
	"log"
	"testing"

	domain "Task-Management/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockSingleResult mocks MongoDB single result methods
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// MockCollection mocks MongoDB collection methods
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) *MockSingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*MockSingleResult)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	cursor, _ := args.Get(0).(*mongo.Cursor) // Ensure it returns *mongo.Cursor
	return cursor, args.Error(1)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// RepositoryTestSuite groups all repository-related tests
type RepositoryTestSuite struct {
	suite.Suite
	client   *mongo.Client
	db       *mongo.Database
	taskRepo domain.TaskRepository
	userRepo domain.UserRepository
}

// SetupSuite runs once before all tests
func (suite *RepositoryTestSuite) SetupSuite() {
	// Connect to the MongoDB test database
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Use a test database
	suite.client = client
	suite.db = client.Database("test_db")

	// Initialize repositories
	suite.taskRepo = NewTaskRepository(suite.db)
	suite.userRepo = NewUserRepository(suite.db)
}

// TearDownSuite runs once after all tests
func (suite *RepositoryTestSuite) TearDownSuite() {
	// Drop the test database
	err := suite.db.Drop(context.Background())
	if err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Clear the users collection before each test
	err := suite.db.Collection(domain.UserCollection).Drop(context.Background())
	if err != nil {
		log.Fatalf("Failed to clear users collection: %v", err)
	}
}

// TaskRepository Tests
func (suite *RepositoryTestSuite) TestTaskRepository_Create() {
	mockTask := &domain.Task{
		Title:  "Test Task",
		UserID: primitive.NewObjectID(),
	}

	result, err := suite.taskRepo.Create(context.Background(), mockTask)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result.ID)
	assert.Equal(suite.T(), "Test Task", result.Title)
}

func (suite *RepositoryTestSuite) TestTaskRepository_GetByID() {
	mockTask := &domain.Task{
		Title:  "Test Task",
		UserID: primitive.NewObjectID(),
	}

	createdTask, err := suite.taskRepo.Create(context.Background(), mockTask)
	assert.NoError(suite.T(), err)

	result, err := suite.taskRepo.GetByID(context.Background(), createdTask.ID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdTask.ID, result.ID)
	assert.Equal(suite.T(), "Test Task", result.Title)
}

func (suite *RepositoryTestSuite) TestTaskRepository_Delete() {
	mockTask := &domain.Task{
		Title:  "Test Task",
		UserID: primitive.NewObjectID(),
	}

	createdTask, err := suite.taskRepo.Create(context.Background(), mockTask)
	assert.NoError(suite.T(), err)

	err = suite.taskRepo.Delete(context.Background(), createdTask.ID)
	assert.NoError(suite.T(), err)

	result, err := suite.taskRepo.GetByID(context.Background(), createdTask.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *RepositoryTestSuite) TestTaskRepository_GetByUserID() {
	mockUserID := primitive.NewObjectID()
	mockTask1 := &domain.Task{Title: "Task 1", UserID: mockUserID}
	mockTask2 := &domain.Task{Title: "Task 2", UserID: mockUserID}

	_, err := suite.taskRepo.Create(context.Background(), mockTask1)
	assert.NoError(suite.T(), err)
	_, err = suite.taskRepo.Create(context.Background(), mockTask2)
	assert.NoError(suite.T(), err)

	tasks, err := suite.taskRepo.GetByUserID(context.Background(), mockUserID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), tasks, 2)
}

func (suite *RepositoryTestSuite) TestTaskRepository_GetAll() {
	mockTask1 := &domain.Task{Title: "Task 1", UserID: primitive.NewObjectID()}
	mockTask2 := &domain.Task{Title: "Task 2", UserID: primitive.NewObjectID()}

	_, err := suite.taskRepo.Create(context.Background(), mockTask1)
	assert.NoError(suite.T(), err)
	_, err = suite.taskRepo.Create(context.Background(), mockTask2)
	assert.NoError(suite.T(), err)

	tasks, err := suite.taskRepo.GetAll(context.Background())
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(tasks), 2)
}

func (suite *RepositoryTestSuite) TestTaskRepository_Update() {
	mockTask := &domain.Task{Title: "Original Title", UserID: primitive.NewObjectID()}
	createdTask, err := suite.taskRepo.Create(context.Background(), mockTask)
	assert.NoError(suite.T(), err)

	createdTask.Title = "Updated Title"
	err = suite.taskRepo.Update(context.Background(), createdTask)
	assert.NoError(suite.T(), err)

	updatedTask, err := suite.taskRepo.GetByID(context.Background(), createdTask.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Title", updatedTask.Title)
}

// UserRepository Tests
func (suite *RepositoryTestSuite) TestUserRepository_Create() {
	mockUser := &domain.User{
		Email: "test@example.com",
	}

	result, err := suite.userRepo.Create(context.Background(), mockUser)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result.ID)
	assert.Equal(suite.T(), "test@example.com", result.Email)
}

func (suite *RepositoryTestSuite) TestUserRepository_GetByEmail() {
	mockUser := &domain.User{
		Email: "test@example.com",
	}

	// Create the user and get the returned user with the generated ID
	createdUser, err := suite.userRepo.Create(context.Background(), mockUser)
	assert.NoError(suite.T(), err)

	// Log the created user ID
	suite.T().Logf("Created User ID: %v", createdUser.ID)

	// Fetch the user by email
	result, err := suite.userRepo.GetByEmail(context.Background(), "test@example.com")

	// Log the expected and actual IDs for debugging
	suite.T().Logf("Expected ID: %v", createdUser.ID)
	suite.T().Logf("Actual ID: %v", result.ID)

	// Assert that the fetched user matches the created user
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdUser.ID, result.ID) // Use the ID returned by Create
	assert.Equal(suite.T(), "test@example.com", result.Email)
}

func (suite *RepositoryTestSuite) TestUserRepository_Delete() {
	mockUser := &domain.User{
		Email: "test@example.com",
	}

	createdUser, err := suite.userRepo.Create(context.Background(), mockUser)
	assert.NoError(suite.T(), err)

	err = suite.userRepo.Delete(context.Background(), createdUser.ID)
	assert.NoError(suite.T(), err)

	result, err := suite.userRepo.GetByID(context.Background(), createdUser.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *RepositoryTestSuite) TestUserRepository_GetByID_NotFound() {
	nonExistentID := primitive.NewObjectID()

	result, err := suite.userRepo.GetByID(context.Background(), nonExistentID)

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *RepositoryTestSuite) TestUserRepository_GetByEmail_NotFound() {
	result, err := suite.userRepo.GetByEmail(context.Background(), "nonexistent@example.com")

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *RepositoryTestSuite) TestUserRepository_Update_Error() {
	mockUser := &domain.User{
		ID: primitive.NewObjectID(),
	}

	err := suite.userRepo.Update(context.Background(), mockUser)

	assert.Error(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestUserRepository_GetAll() {
	mockUser1 := &domain.User{Email: "user1@example.com"}
	mockUser2 := &domain.User{Email: "user2@example.com"}

	_, err := suite.userRepo.Create(context.Background(), mockUser1)
	assert.NoError(suite.T(), err)
	_, err = suite.userRepo.Create(context.Background(), mockUser2)
	assert.NoError(suite.T(), err)

	users, err := suite.userRepo.GetAll(context.Background())
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(users), 2)
}

// Run the test suite
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
