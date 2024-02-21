package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		log.Fatal(err)
	}
	mongodb := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb))
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error while connecting mongo client")
	}
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// client, err := mongo.NewClient(options.Client().ApplyURI(mongodb))

	// if err != nil {
	// 	fmt.Println("Found error while connecting mongo db:", err.Error())
	// 	log.Fatal(err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// err = client.Conect(ctx)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("Connected to mongoDB")

	return client

}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster").Collection(collectionName)
	fmt.Println("Collection Name :", collection)
	fmt.Println("Collection Type :", reflect.TypeOf(collection))

	return collection

}
