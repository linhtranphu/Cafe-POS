package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/printing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type printNotificationRepository struct {
	collection *mongo.Collection
}

// NewPrintNotificationRepository creates a new print notification repository
func NewPrintNotificationRepository(db *mongo.Database) printing.PrintNotificationRepository {
	return &printNotificationRepository{
		collection: db.Collection("print_notifications"),
	}
}

// Create creates a new notification
func (r *printNotificationRepository) Create(notification *printing.PrintNotification) error {
	if notification.ID.IsZero() {
		notification.ID = primitive.NewObjectID()
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

// FindUnread returns all unread notifications
func (r *printNotificationRepository) FindUnread() ([]*printing.PrintNotification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"read": false}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*printing.PrintNotification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

// FindByJobID returns all notifications for a specific job
func (r *printNotificationRepository) FindByJobID(jobID primitive.ObjectID) ([]*printing.PrintNotification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"job_id": jobID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*printing.PrintNotification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (r *printNotificationRepository) MarkAsRead(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"read": true}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// MarkAllAsRead marks all notifications as read
func (r *printNotificationRepository) MarkAllAsRead() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"read": false}
	update := bson.M{"$set": bson.M{"read": true}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// DeleteOld deletes notifications older than the specified time
func (r *printNotificationRepository) DeleteOld(olderThan time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"created_at": bson.M{"$lt": olderThan}}

	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
