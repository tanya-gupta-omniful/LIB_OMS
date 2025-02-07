package service

import (
	"OMS/domain"
	"OMS/repo"
	"OMS/utils/sqs"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/omniful/go_commons/log"
)

type Service interface {
	ProcessOrder(filePath string) error
	CreateBulkOrderService(filePath string) error
}

type service struct {
	repo repo.Repository
}

// NewService is the constructor function to create a new instance of ConcreteService.
func NewService(r repo.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) ProcessOrder(filePath string) error {
	// Check if the file exists
	_, err := os.Stat(filePath)

	fmt.Println(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist")
		}
		return fmt.Errorf("unable to access file: %v", err)
	}

	bulkOrderEvent := domain.BulkOrderEvent{
		FilePath: filePath,
		User: domain.User{
			ID:        "4",
			FirstName: "Amaan",
			LastName:  "Goyal",
			Email:     "ghhh@yahoo",
		},
		RequestTime: time.Now(),
	}

	err = sqs.PushEmailMessageToSQS(context.Background(), bulkOrderEvent)
	if err != nil {
		log.Errorf("sqs push error %w", err)
	}
	// Open the CSV file
//file, err := os.Open(filePath)
//if err != nil {
//	return fmt.Errorf("unable to open file: %v", err)
//}
//defer file.Close()
//

//// Read and parse the CSV content
//reader := csv.NewReader(file)
//records, err := reader.ReadAll()
//if err != nil {
//	return fmt.Errorf("unable to read CSV content: %v", err)
//}

//for _, record := range records {
//	// Save the record to MongoDB via the repository
//	err := s.repo.SaveOrder(record)
//	if err != nil {
//		return fmt.Errorf("unable to save record to MongoDB: %v", err)
//	}
//}

	return nil
}

func (s *service) CreateBulkOrderService(filePath string) error {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist")
		}
		return fmt.Errorf("unable to access file: %v", err)
	}

	// Assuming the event has user info and timestamp, you may adjust accordingly
	bulkOrderEvent := domain.BulkOrderEvent{
		FilePath: filePath,
		User: domain.User{
			ID:        "4",
			FirstName: "Amaan",
			LastName:  "Goyal",
			Email:     "ghhh@yahoo",
		},
		RequestTime: time.Now(),
	}

	// Pass the BulkOrderEvent to the repository to process and save orders
	err = s.repo.SaveBulkOrders(bulkOrderEvent)
	if err != nil {
		log.Errorf("Failed to process bulk order: %v", err)
		return err
	}

	return nil
}