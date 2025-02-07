package appinit

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
var (
	mongoClient     *mongo.Client
	OrderCollection *mongo.Collection
)

func Initialize(ctx context.Context) {
	initializeDB(ctx)
	initializeLog(ctx)
}
func initializeLog(ctx context.Context) {
	err := log.InitializeLogger(
		log.Formatter(config.GetString(ctx, "log.format")),
		log.Level(config.GetString(ctx, "log.level")),
	)
	if err != nil {
		log.WithError(err).Panic("unable to initialise log")
	}
}
func initializeDB(ctx context.Context){
	dbName := config.GetString(ctx, "mongo.database")
	collectionName := config.GetString(ctx, "mongo.collection")
	var err error
	
	mongoURI :=config.GetString(ctx, "mongo.uri")
	if mongoURI == "" {
		log.Panic("MongoDB URI is not set in configuration")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil{
		fmt.Println("Error connecting to db",err)
		return
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err!=nil{
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}
    
	fmt.Println("Mongodb cnnection succesfull")
	OrderCollection = mongoClient.Database(dbName).Collection(collectionName)
}