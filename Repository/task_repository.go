package repository

import (
	"context"
	"errors"
	"time"

	domain "Task-Management/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CollectionInterface abstracts MongoDB collection operations
type CollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}

// MongoCollectionWrapper wraps *mongo.Collection to implement CollectionInterface
type MongoCollectionWrapper struct {
	collection *mongo.Collection
}

func (m *MongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(ctx, document)
}

func (m *MongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	return m.collection.FindOne(ctx, filter)
}

func (m *MongoCollectionWrapper) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return m.collection.Find(ctx, filter)
}

func (m *MongoCollectionWrapper) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return m.collection.DeleteOne(ctx, filter)
}

func (m *MongoCollectionWrapper) UpdateOne(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error) {
	return m.collection.UpdateOne(ctx, filter, update)
}

// TaskRepository defines the expected behavior for the task repository
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Task, error)
	GetAll(ctx context.Context) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type taskRepository struct {
	collection CollectionInterface
}

// NewTaskRepository initializes a new task repository
func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &taskRepository{
		collection: &MongoCollectionWrapper{collection: db.Collection(domain.TaskCollection)},
	}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to parse inserted ID as ObjectID")
	}
	task.ID = id
	return task, nil
}

func (r *taskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var task domain.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no document is found
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	var tasks []*domain.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	var tasks []*domain.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found to update")
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
