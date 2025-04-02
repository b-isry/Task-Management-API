package Usecases

import (
	domain "Task-Management/Domain"
	infrastructure "Task-Management/Infrastructure"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUseCase struct {
	userRepo         domain.UserRepository
	hashPassword     func(string) (string, error)
	comparePasswords func(string, string) bool
	generateToken    func(string, string) (string, error)
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo:         userRepo,
		hashPassword:     infrastructure.HashPassword,     // Default implementation
		comparePasswords: infrastructure.ComparePasswords, // Default implementation
		generateToken:    infrastructure.GenerateToken,    // Default implementation
	}
}

func (u *userUseCase) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && err.Error() != "user not found" { // Adjust error check
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return u.userRepo.Create(ctx, user)
}

func (u *userUseCase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !u.comparePasswords(user.Password, password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := u.generateToken(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *userUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *userUseCase) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.Password != "" {
		hashedPassword, err := u.hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	return u.userRepo.Update(ctx, user)
}

func (u *userUseCase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return u.userRepo.Delete(ctx, id)
}
