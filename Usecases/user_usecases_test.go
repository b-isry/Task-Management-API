package Usecases

import (
	"context"
	"errors"
	"testing"

	"Task-Management/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*Domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *Domain.User) (*Domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*Domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *Domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetUserByEmail retrieves a user by email
func (u *userUseCase) GetUserByEmail(ctx context.Context, email string) (*Domain.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

// CreateUser creates a new user
func (u *userUseCase) CreateUser(ctx context.Context, user *Domain.User) (*Domain.User, error) {
	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	return u.userRepo.Create(ctx, user)
}

// UserUseCaseTestSuite groups all user use case-related tests
type UserUseCaseTestSuite struct {
	suite.Suite
	mockRepo     *MockUserRepository
	mockHashFunc func(string) (string, error)
	userUseCase  *userUseCase
}

// SetupSuite runs once before all tests
func (suite *UserUseCaseTestSuite) SetupSuite() {
	suite.mockHashFunc = func(password string) (string, error) {
		return "hashedPassword", nil
	}
}

// SetupTest runs before each test
func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockUserRepository)
	suite.userUseCase = &userUseCase{
		userRepo:         suite.mockRepo,
		hashPassword:     suite.mockHashFunc,
		comparePasswords: func(hashedPassword, plainPassword string) bool { return true },
		generateToken:    func(userID, role string) (string, error) { return "mockToken", nil },
	}
}

// TestUpdateUser tests updating a user successfully
func (suite *UserUseCaseTestSuite) TestUpdateUser() {
	mockUser := &Domain.User{
		ID:       primitive.NewObjectID(),
		Email:    "test@example.com",
		Password: "newPassword",
	}

	// Mock repository behavior
	suite.mockRepo.On("Update", mock.Anything, mockUser).Return(nil)

	err := suite.userUseCase.UpdateUser(context.Background(), mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "hashedPassword", mockUser.Password)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUpdateUser_RepositoryError tests repository error during user update
func (suite *UserUseCaseTestSuite) TestUpdateUser_RepositoryError() {
	mockUser := &Domain.User{
		ID:       primitive.NewObjectID(),
		Email:    "test@example.com",
		Password: "newPassword",
	}

	// Mock repository behavior
	suite.mockRepo.On("Update", mock.Anything, mockUser).Return(errors.New("repository error"))

	err := suite.userUseCase.UpdateUser(context.Background(), mockUser)

	assert.EqualError(suite.T(), err, "repository error")
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestDeleteUser tests deleting a user successfully
func (suite *UserUseCaseTestSuite) TestDeleteUser() {
	userID := primitive.NewObjectID()

	// Mock repository behavior
	suite.mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	err := suite.userUseCase.DeleteUser(context.Background(), userID)

	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetUserByEmail tests fetching a user by email successfully
func (suite *UserUseCaseTestSuite) TestGetUserByEmail() {
	email := "test@example.com"
	mockUser := &Domain.User{
		ID:    primitive.NewObjectID(),
		Email: email,
	}

	// Mock repository behavior
	suite.mockRepo.On("GetByEmail", mock.Anything, email).Return(mockUser, nil)

	result, err := suite.userUseCase.GetUserByEmail(context.Background(), email)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), email, result.Email)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateUser tests creating a user successfully
func (suite *UserUseCaseTestSuite) TestCreateUser() {
	mockUser := &Domain.User{
		Email:    "newuser@example.com",
		Password: "password123",
	}

	// Mock repository behavior
	suite.mockRepo.On("Create", mock.Anything, mockUser).Return(mockUser, nil)

	result, err := suite.userUseCase.CreateUser(context.Background(), mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "hashedPassword", result.Password)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetAllUsers tests fetching all users successfully
func (suite *UserUseCaseTestSuite) TestGetAllUsers() {
	mockUsers := []*Domain.User{
		{ID: primitive.NewObjectID(), Email: "user1@example.com"},
		{ID: primitive.NewObjectID(), Email: "user2@example.com"},
	}

	// Mock repository behavior
	suite.mockRepo.On("GetAll", mock.Anything).Return(mockUsers, nil)

	results, err := suite.userUseCase.GetAllUsers(context.Background())

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetUserByID tests fetching a user by ID successfully
func (suite *UserUseCaseTestSuite) TestGetUserByID() {
	userID := primitive.NewObjectID()
	mockUser := &Domain.User{ID: userID, Email: "user@example.com"}

	// Mock repository behavior
	suite.mockRepo.On("GetByID", mock.Anything, userID).Return(mockUser, nil)

	result, err := suite.userUseCase.GetUserByID(context.Background(), userID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, result.ID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetUserByEmail_NotFound tests fetching a user by email when not found
func (suite *UserUseCaseTestSuite) TestGetUserByEmail_NotFound() {
	email := "nonexistent@example.com"

	// Mock repository behavior
	suite.mockRepo.On("GetByEmail", mock.Anything, email).Return(nil, errors.New("user not found"))

	result, err := suite.userUseCase.GetUserByEmail(context.Background(), email)

	assert.Nil(suite.T(), result)
	assert.EqualError(suite.T(), err, "user not found")
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestRegisterUser tests registering a user successfully
func (suite *UserUseCaseTestSuite) TestRegisterUser() {
	mockUser := &Domain.User{
		Email:    "newuser@example.com",
		Password: "password123",
	}

	// Mock repository behavior
	suite.mockRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, errors.New("user not found"))
	suite.mockRepo.On("Create", mock.Anything, mockUser).Return(mockUser, nil)

	result, err := suite.userUseCase.Register(context.Background(), mockUser)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "hashedPassword", result.Password)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestRegisterUser_UserAlreadyExists tests registering a user when the user already exists
func (suite *UserUseCaseTestSuite) TestRegisterUser_UserAlreadyExists() {
	mockUser := &Domain.User{
		Email:    "existinguser@example.com",
		Password: "password123",
	}

	// Mock repository behavior
	suite.mockRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil)

	result, err := suite.userUseCase.Register(context.Background(), mockUser)

	assert.Nil(suite.T(), result)
	assert.EqualError(suite.T(), err, "user already exists")
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestLoginUser tests logging in a user successfully
func (suite *UserUseCaseTestSuite) TestLoginUser() {
	email := "user@example.com"
	password := "password123"
	mockUser := &Domain.User{
		ID:       primitive.NewObjectID(),
		Email:    email,
		Password: "hashedPassword",
		Role:     "user",
	}

	// Mock repository behavior
	suite.mockRepo.On("GetByEmail", mock.Anything, email).Return(mockUser, nil)

	result, token, err := suite.userUseCase.Login(context.Background(), email, password)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), email, result.Email)
	assert.Equal(suite.T(), "mockToken", token)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestLoginUser_InvalidCredentials tests logging in with invalid credentials
func (suite *UserUseCaseTestSuite) TestLoginUser_InvalidCredentials() {
	email := "user@example.com"
	password := "wrongPassword"
	mockUser := &Domain.User{
		ID:       primitive.NewObjectID(),
		Email:    email,
		Password: "hashedPassword",
	}

	// Mock repository behavior
	suite.mockRepo.On("GetByEmail", mock.Anything, email).Return(mockUser, nil)
	suite.userUseCase.comparePasswords = func(hashedPassword, plainPassword string) bool { return false }

	result, token, err := suite.userUseCase.Login(context.Background(), email, password)

	assert.Nil(suite.T(), result)
	assert.Empty(suite.T(), token)
	assert.EqualError(suite.T(), err, "invalid credentials")
	suite.mockRepo.AssertExpectations(suite.T())
}

// Run the test suite
func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}
