package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        string `json:"id"`        // User ID
	FirstName string `json:"firstName"` // User's first name
	LastName  string `json:"lastName"`  // User's last name
	Email     string `json:"email"`     // User's email address
}

type BulkOrderEvent struct {
	FilePath    string    `json:"filePath"`    // Path to the uploaded file
	User        User      `json:"user"`        // User who triggered the request
	RequestTime time.Time `json:"requestTime"` // Timestamp of the request
}
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OrderID     string             `bson:"order_id"`     
	TenantID    string             `bson:"tenant_id"`    
	SellerID    string             `bson:"seller_id"`    
	HubID       string             `bson:"hub_id"`       
	CustomerID  string             `bson:"customer_id"`  
	Items       []OrderItem        `bson:"items"`        
	TotalAmount float64            `bson:"total_amount"` 
	Status      string             `bson:"status"`       
	CreatedAt   time.Time          `bson:"created_at"`   
	UpdatedAt   time.Time          `bson:"updated_at"`   
}
type OrderItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	SKUID   string  `bson:"sku_id"`   
	Quantity int     `bson:"qtty"`    
	Price    float64 `bson:"price"`   
}