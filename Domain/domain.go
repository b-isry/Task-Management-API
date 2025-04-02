package Domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserCollection = "users"
	RoleUser       = "user"
	RoleAdmin      = "admin"
)

const (
	TaskCollection   = "tasks"
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

// User represents the core user entity
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Task represents the core task entity
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	Status      string             `bson:"status" json:"status"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*Task, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Task, error)
	GetAll(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	Register(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, email, password string) (*User, string, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}

// TaskUseCase defines the interface for task business logic
type TaskUseCase interface {
	CreateTask(ctx context.Context, task *Task) (*Task, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Task, error)
	GetTasksByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Task, error)
	GetAllTasks(ctx context.Context) ([]*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
}

// Request/Response DTOs
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

