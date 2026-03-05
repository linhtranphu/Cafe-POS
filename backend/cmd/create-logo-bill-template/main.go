package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")
	collection := db.Collection("print_templates")

	// New template with logo and table format
	templateContent := `{{if .ShowLogo}}[LOGO]

{{end}}{{truncate .ShopName 48}}
{{if .ShowAddress}}{{if .ShopAddress}}{{truncate .ShopAddress 48}}
{{end}}{{end}}{{if .ShowPhone}}{{if .ShopPhone}}Tel: {{.ShopPhone}}
{{end}}{{end}}================================
HÓA ĐƠN BÁN HÀNG
================================
Order: {{.Order.OrderNumber}}
Ngày: {{formatTime .Order.CreatedAt "02/01/2006 15:04"}}
{{if .Order.WaiterName}}Phục vụ: {{truncate .Order.WaiterName 40}}
{{end}}================================
[TABLE_START]
Tên món              SL  Đơn giá    Thành tiền
------------------------------------------------
{{range .Order.Items}}{{truncate .Name 48}}{{if .VariantName}}
  ({{truncate .VariantName 44}}){{end}}
  {{.Quantity}}  {{formatPrice .Price}}  {{formatPrice .Subtotal}}
{{end}}[TABLE_END]
================================
Tổng tiền: {{formatPrice .Order.Subtotal}} VND
{{if gt .Order.Discount 0.0}}Giảm giá: -{{formatPrice .Order.Discount}} VND
{{end}}--------------------------------
TỔNG CỘNG: {{formatPrice .Order.Total}} VND
================================
{{if .ShowCustomMessage}}{{if .CustomMessage}}{{truncate .CustomMessage 48}}
{{end}}{{end}}Cảm ơn quý khách!
Hẹn gặp lại!
`

	template := bson.M{
		"_id":         primitive.NewObjectID(),
		"name":        "Bill với Logo và Bảng",
		"type":        "BILL",
		"content":     templateContent,
		"description": "Template hóa đơn với logo ở góc trên bên trái và bảng món có cấu trúc rõ ràng",
		"is_default":  false, // Don't override existing default
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	result, err := collection.InsertOne(context.TODO(), template)
	if err != nil {
		log.Fatalf("Failed to insert template: %v", err)
	}

	fmt.Printf("✅ Template created successfully with ID: %s\n", result.InsertedID)
	fmt.Println("📝 Template name: Bill với Logo và Bảng")
	fmt.Println("🔧 To set as default, go to Print Management > Templates")
}
