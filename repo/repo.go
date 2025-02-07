package repo

import (
	"OMS/domain"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	//"github.com/omniful/go_commons/db/sql/postgres"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	SaveOrder(record []string) error
	SaveBulkOrders(orders []domain.Order) error
}

// Inter
type repository struct {
	db *mongo.Client
}

// Singleton instance of Repository and the sync.Once to ensure it's initialized only once.
var repo *repository
var repoOnce sync.Once

// NewRepository is the constructor function that ensures the Repository is initialized only once.
func NewRepository(db *mongo.Client) Repository {
	repoOnce.Do(func() {
		// Initialize the Repository with a given DbCluster.
		repo = &repository{
			db: db,
		}
	})
	return repo
}

func (r *repository) SaveOrder(record []string) error {
	collection := r.db.Database("oms").Collection("orders")

	// Create a document from the CSV record
	document := bson.D{
		{"column1", record[0]},
		{"column2", record[1]},
		{"column3", record[2]},
		// Add more fields as needed based on your CSV structure
	}

	// Insert the document into MongoDB
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	return nil
}

func (r *repository) SaveBulkOrders(bulkOrderEvent domain.BulkOrderEvent) error {

	file, err := os.Open(bulkOrderEvent.FilePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read CSV content: %v", err)
	}

	var validOrders []domain.Order

	for _, record := range records {
		isValid, err := r.validateOrder(record[0], record[1], record[2]) // SKU, HubID, Quantity
		if err != nil {
			log.Printf("Skipping invalid order %v due to error: %v", record, err)
			continue
		}
		if isValid {
			// Create an OrderItem and Order (assumes validation passed)
			order := domain.Order{
				OrderID:   record[0], // Example, set properly
				HubID:     record[1], // Example, set properly
				Items: []domain.OrderItem{
					{
						SKUID:   record[0],  // SKU from the CSV
						Quantity: record[2], // Quantity from the CSV
					},
				},
			}
			validOrders = append(validOrders, order)
		}
	}

	// Save valid orders to MongoDB in bulk
	if len(validOrders) > 0 {
		collection := r.db.Database("orders").Collection("orders")
		var documents []interface{}
		for _, order := range validOrders {
			documents = append(documents, bson.D{
				{"order_id", order.OrderID},
				{"hub_id", order.HubID},
				{"items", order.Items},
			})
		}
		_, err = collection.InsertMany(context.Background(), documents)
		if err != nil {
			return fmt.Errorf("failed to insert bulk orders: %v", err)
		}
	}

	return nil
}

func (r *repository) validateOrder(skuID, hubID, quantity string) (bool, error) {
	// Make the API call to validate the SKU, Hub, and Quantity
	apiURL := "http://localhost:8120/api/v1/inventory/validate" 

	// Build the request parameters or payload
	payload := fmt.Sprintf("sku_id=%s&hub_id=%s&qtty=%s", skuID, hubID, quantity)

	// Call the API and handle response
	response, err := http.PostForm(apiURL, payload)
	if err != nil {
		return false, fmt.Errorf("failed to call validation API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return false, fmt.Errorf("validation API returned non-200 status code: %v", response.StatusCode)
	}

	var result struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	if result.Success {
		// If validation is successful, call InsertOrder to save it to DB
		err := r.InsertOrder(skuID, hubID, quantity)
		if err != nil {
			return false, fmt.Errorf("failed to insert order: %v", err)
		}
		return true, nil
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to parse API response: %v", err)
	}

	return result.Success, nil
}
func (r *repository) InsertOrder(skuID, hubID, quantity string) error {
	// Create an OrderItem from the validated data
	orderItem := domain.OrderItem{
		SKUID:   skuID,    // SKU from the validated data
		Quantity: quantity, // Quantity from the validated data
	}

	// Create an Order (you can set other fields like OrderID, CustomerID, etc.)
	order := domain.Order{
		OrderID:   generateOrderID(), // You may need a function to generate a unique Order ID
		HubID:     hubID,
		Items:     []domain.OrderItem{orderItem},
		Status:    "pending", // Set the initial order status
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Access the database and collection
	collection := r.db.Database("orders").Collection("orders")

	// Insert the new order into the database
	_, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		return fmt.Errorf("failed to insert order into database: %v", err)
	}

	return nil
}