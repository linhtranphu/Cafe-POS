package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AICommandLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Message    string             `bson:"message" json:"message"`
	Role       string             `bson:"role" json:"role"` // "user" or "agent"
	ActionType string             `bson:"action_type,omitempty" json:"action_type,omitempty"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}

type AILogRepository struct {
	collection *mongo.Collection
}

func NewAILogRepository(db *mongo.Database) *AILogRepository {
	return &AILogRepository{
		collection: db.Collection("ai_command_logs"),
	}
}

func (r *AILogRepository) Insert(ctx context.Context, log *AICommandLog) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	log.ID = primitive.NewObjectID()
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	_, err := r.collection.InsertOne(ctx, log)
	return err
}

func (r *AILogRepository) GetRecent(ctx context.Context, limit int) ([]AICommandLog, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []AICommandLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
