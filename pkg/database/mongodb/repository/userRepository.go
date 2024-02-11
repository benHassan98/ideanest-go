package repository

import (
	"Ideanest/pkg/database/mongodb/models"
	"Ideanest/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	config     utils.Config
	client     *mongo.Client
	collection string
}

func NewUserRepository(config utils.Config, client *mongo.Client) UserRepository {
	return UserRepository{config: config, client: client, collection: "users"}
}

func (u UserRepository) InsertOne(ctx context.Context, user models.User) (*mongo.InsertOneResult, error) {

	collection := u.client.Database(u.config.Database.DbName).Collection(u.collection)

	insertOneResult, err := collection.InsertOne(ctx, user)

	return insertOneResult, err
}

func (u UserRepository) FindOne(ctx context.Context, email string, password string) error {

	collection := u.client.Database(u.config.Database.DbName).Collection(u.collection)

	err := collection.FindOne(ctx, bson.D{{"email", email}, {"password", password}}).Err()

	return err
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {

	var user models.User

	collection := u.client.Database(u.config.Database.DbName).Collection(u.collection)

	err := collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)

	if err != nil {
		return user, err
	}

	return user, nil
}
