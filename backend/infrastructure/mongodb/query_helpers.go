package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QueryTimeout is the default timeout for MongoDB queries
const QueryTimeout = 5 * time.Second

// WithQueryTimeout wraps a context with a timeout for MongoDB queries
func WithQueryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, QueryTimeout)
}

// IsCollectionNotFoundError checks if the error is due to collection not existing
func IsCollectionNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for common "collection not found" errors
	if err == mongo.ErrNoDocuments {
		return true
	}
	
	// Check if it's a namespace not found error
	if mongo.IsDuplicateKeyError(err) {
		return false
	}
	
	// MongoDB returns this when collection doesn't exist
	errMsg := err.Error()
	return errMsg == "ns not found" || 
		   errMsg == "collection does not exist" ||
		   errMsg == "namespace not found"
}

// SafeFindAll executes a Find query and returns empty slice if collection doesn't exist
// This prevents timeout errors when querying non-existent collections
func SafeFindAll[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]*T, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		// Return empty slice if collection doesn't exist
		if IsCollectionNotFoundError(err) {
			return []*T{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*T
	if err = cursor.All(ctx, &results); err != nil {
		// Return empty slice if collection doesn't exist
		if IsCollectionNotFoundError(err) {
			return []*T{}, nil
		}
		return nil, err
	}
	
	// Ensure we never return nil slice
	if results == nil {
		results = []*T{}
	}
	
	return results, nil
}

// SafeFindOne executes a FindOne query with timeout
func SafeFindOne[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.FindOneOptions) (*T, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	var result T
	err := collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// SafeCount executes a CountDocuments query with timeout
func SafeCount(ctx context.Context, collection *mongo.Collection, filter interface{}) (int64, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		// Return 0 if collection doesn't exist
		if IsCollectionNotFoundError(err) {
			return 0, nil
		}
		return 0, err
	}
	
	return count, nil
}

// SafeInsertOne executes an InsertOne with timeout
func SafeInsertOne(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	return collection.InsertOne(ctx, document)
}

// SafeUpdateOne executes an UpdateOne with timeout
func SafeUpdateOne(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	return collection.UpdateOne(ctx, filter, update)
}

// SafeDeleteOne executes a DeleteOne with timeout
func SafeDeleteOne(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	return collection.DeleteOne(ctx, filter)
}

// SafeDeleteMany executes a DeleteMany with timeout
func SafeDeleteMany(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	// Add timeout
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	return collection.DeleteMany(ctx, filter)
}
