package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectToMongo(ctx context.Context, url string) (*mongo.Client, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Println("There was a problem connecting to your Atlas cluster. Check that the URI includes a valid username and password, and that your IP address has been added to the access list. Error: ")
		return nil, err
	}

	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}
