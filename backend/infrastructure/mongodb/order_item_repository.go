package mongodb

import (
	"context"
	"time"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemRepository struct {
	collection *mongo.Collection
}

func NewOrderItemRepository(db *mongo.Database) *OrderItemRepository {
	return &OrderItemRepository{
		collection: db.Collection("order_items"),
	}
}

// CreateIndexes creates the required indexes for the order_items collection
func (r *OrderItemRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"order_id", 1}},
			Options: options.Index().SetName("idx_order_id"),
		},
		{
			Keys: bson.D{{"menu_item_id", 1}},
			Options: options.Index().SetName("idx_menu_item_id"),
		},
		{
			Keys: bson.D{{"cost_status", 1}},
			Options: options.Index().SetName("idx_cost_status"),
		},
		{
			Keys: bson.D{{"cost_calculated_at", 1}},
			Options: options.Index().SetName("idx_cost_calculated_at"),
		},
		{
			Keys: bson.D{{"order_id", 1}, {"menu_item_id", 1}},
			Options: options.Index().SetName("idx_order_menu_item"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create inserts a new order item with cost tracking
func (r *OrderItemRepository) Create(ctx context.Context, item *order.OrderItemWithCost) error {
	item.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// CreateMany inserts multiple order items with cost tracking
func (r *OrderItemRepository) CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error {
	if len(items) == 0 {
		return nil
	}

	docs := make([]interface{}, len(items))
	for i, item := range items {
		item.CreatedAt = time.Now()
		docs[i] = item
	}

	results, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	// Update IDs
	for i, id := range results.InsertedIDs {
		items[i].ID = id.(primitive.ObjectID)
	}

	return nil
}

// FindByOrderID retrieves all order items for a specific order
func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*order.OrderItemWithCost
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// FindByMenuItemID retrieves all order items for a specific menu item
func (r *OrderItemRepository) FindByMenuItemID(ctx context.Context, menuItemID primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"menu_item_id": menuItemID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*order.OrderItemWithCost
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// FindByCostStatus retrieves all order items with a specific cost status
func (r *OrderItemRepository) FindByCostStatus(ctx context.Context, status order.CostStatus) ([]*order.OrderItemWithCost, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"cost_status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*order.OrderItemWithCost
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateCost updates the accounting cost for an order item
func (r *OrderItemRepository) UpdateCost(ctx context.Context, id primitive.ObjectID, accountingCost float64, costStatus order.CostStatus) error {
	update := bson.M{
		"$set": bson.M{
			"accounting_cost":    accountingCost,
			"cost_calculated_at": time.Now(),
			"cost_status":        costStatus,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateManyCostsByOrderID updates the accounting cost for all items in an order
func (r *OrderItemRepository) UpdateManyCostsByOrderID(ctx context.Context, orderID primitive.ObjectID, accountingCost float64, costStatus order.CostStatus) error {
	update := bson.M{
		"$set": bson.M{
			"accounting_cost":    accountingCost,
			"cost_calculated_at": time.Now(),
			"cost_status":        costStatus,
		},
	}
	_, err := r.collection.UpdateMany(ctx, bson.M{"order_id": orderID}, update)
	return err
}

// FindByDateRange retrieves order items within a date range
func (r *OrderItemRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.OrderItemWithCost, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*order.OrderItemWithCost
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// FindByOrderIDs retrieves order items for multiple orders
func (r *OrderItemRepository) FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	filter := bson.M{
		"order_id": bson.M{"$in": orderIDs},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*order.OrderItemWithCost
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteByOrderID deletes all order items for a specific order
func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"order_id": orderID})
	return err
}
