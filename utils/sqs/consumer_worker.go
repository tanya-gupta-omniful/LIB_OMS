package sqs

import (
	psqs "OMS/sqs"
	"context"
	"fmt"
	"log"
	"OMS/service"

	"github.com/omniful/go_commons/sqs"
)

type ExampleHandler struct{}

func (h *ExampleHandler) Process(ctx context.Context, message *[]sqs.Message) error {
	//TODO implement me
	//panic("implement me")

	for _, msg := range *message {
		err := h.Handle(&msg)
		if err != nil {

		}
	}
	return nil
}

func (h *ExampleHandler) Handle(msg *sqs.Message) error {
	fmt.Println("Processing message:", string(msg.Value))
	var event struct {
		FilePath string `json:"filePath"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to parse SQS message: %v", err)
		return err
	}

	// Call service to process the bulk order
	err := service.CreateBulkOrderService(event.FilePath)
	if err != nil {
		log.Printf("Failed to process bulk order: %v", err)
		return err
	}

	fmt.Println("Bulk order processing complete")
	

	return nil
}

func StartConsumerWorker(ctx context.Context) {

	// Set up consumer
	handler := &ExampleHandler{}
	consumer, err := sqs.NewConsumer(
		psqs.QueueGlobal,
		1, // Number of workers
		1, // Concurrency per worker
		handler,
		10,    // Max messages count
		30,    // Visibility timeout
		false, // Is async
		false, // Send batch message
	)

	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	consumer.Start(ctx)

	// Let the consumer run for a while
	//time.Sleep(10 * time.Second)

	//consumer.Close()
}