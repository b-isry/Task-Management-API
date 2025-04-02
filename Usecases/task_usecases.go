package Usecases

import (
	"context"
	"errors"
	"time"

	domain "Task-Management/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewTaskUseCase(taskRepo domain.TaskRepository) domain.TaskUseCase {
	return &taskUseCase{
		taskRepo: taskRepo,
	}
}

func (t *taskUseCase) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	// Validate task
	if task.Title == "" {
		return nil, errors.New("task title is required")
	}
	if task.DueDate.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	// Set initial status
	task.Status = domain.StatusPending

	return t.taskRepo.Create(ctx, task)
}

func (t *taskUseCase) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	return t.taskRepo.GetByID(ctx, id)
}

func (t *taskUseCase) GetTasksByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Task, error) {
	return t.taskRepo.GetByUserID(ctx, userID)
}

func (t *taskUseCase) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	// Fetch all tasks from the repository
	tasks, err := t.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *taskUseCase) UpdateTask(ctx context.Context, task *domain.Task) error {
	// Validate task
	if task.Title == "" {
		return errors.New("task title is required")
	}
	if task.DueDate.Before(time.Now()) {
		return errors.New("due date cannot be in the past")
	}

	// Validate status transition
	existingTask, err := t.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}

	// Only allow status transitions from pending to in_progress to completed
	if existingTask.Status == domain.StatusCompleted && task.Status != domain.StatusCompleted {
		return errors.New("cannot change status of completed task")
	}

	return t.taskRepo.Update(ctx, task)
}

func (t *taskUseCase) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	return t.taskRepo.Delete(ctx, id)
}
