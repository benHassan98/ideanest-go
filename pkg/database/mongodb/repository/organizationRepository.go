package repository

import (
	"Ideanest/pkg/database/mongodb/models"
	"Ideanest/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrganizationRepository struct {
	config     utils.Config
	client     *mongo.Client
	collection string
}

func NewOrganizationRepository(config utils.Config, client *mongo.Client) OrganizationRepository {
	return OrganizationRepository{config: config, client: client, collection: "organizations"}
}

func (o OrganizationRepository) InsertOne(ctx context.Context, organization models.Organization) (*mongo.InsertOneResult, error) {

	collection := o.client.Database(o.config.Database.DbName).Collection(o.collection)

	insertOneResult, err := collection.InsertOne(ctx, organization)

	return insertOneResult, err
}

func (o OrganizationRepository) FindById(ctx context.Context, id primitive.ObjectID) (models.Organization, error) {

	var organization models.Organization

	collection := o.client.Database(o.config.Database.DbName).Collection(o.collection)

	err := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&organization)

	return organization, err

}

func (o OrganizationRepository) FindAll(ctx context.Context) ([]models.Organization, error) {

	var organization []models.Organization

	collection := o.client.Database(o.config.Database.DbName).Collection(o.collection)

	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &organization); err != nil {
		return nil, err
	}

	return organization, nil

}

func (o OrganizationRepository) Update(ctx context.Context, organization models.Organization) error {

	collection := o.client.Database(o.config.Database.DbName).Collection(o.collection)

	_, err := collection.UpdateByID(ctx, organization.Id, bson.D{{"$set", organization}}, nil)

	if err != nil {
		return err
	}

	return nil
}

func (o OrganizationRepository) DeleteById(ctx context.Context, id primitive.ObjectID) error {

	collection := o.client.Database(o.config.Database.DbName).Collection(o.collection)

	_, err := collection.DeleteOne(ctx, bson.D{{"_id", id}})

	if err != nil {
		return err
	}

	return nil
}

func (o OrganizationRepository) InviteMember(ctx context.Context, organizationId primitive.ObjectID, email string) error {

	organization, err := o.FindById(ctx, organizationId)

	if err != nil {
		return err
	}

	user, err := UserRepository{o.config, o.client, "users"}.FindByEmail(ctx, email)

	if err != nil {
		return err
	}

	organization.OrganizationMembers = append(organization.OrganizationMembers, models.Member{Name: user.Name, Email: user.Email, AccessLevel: "invite"})

	collection := o.client.Database(o.config.Database.DbName).Collection(o.collection)

	_, err = collection.UpdateByID(ctx, organization.Id, bson.D{{"$set", organization}}, nil)

	if err != nil {
		return err
	}

	return nil
}
