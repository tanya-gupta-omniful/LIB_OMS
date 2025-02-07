package sqs

import (
	"context"
	"github.com/omniful/go_commons/config"
	
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	
)

var QueueGlobal *sqs.Queue

func IntiializeSqs(ctx context.Context) {
	//QueueGlobal.

	queueName := config.GetString(ctx, "sqs.name")
	queue, err := sqs.NewStandardQueue(ctx, queueName, &sqs.Config{
		Account:  config.GetString(ctx, "sqs.account"),
		Endpoint: config.GetString(ctx, "sqs.endpoint"),
		Region:   config.GetString(ctx, "sqs.region"),
	})

	if err != nil || queue == nil {
		//cusErr = error2.NewCustomError(error3.SqsInitializeErr, err.Error())
		log.Errorf("error init sqs %w", err)
		return
	}

	QueueGlobal = queue

}