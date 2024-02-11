package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty"`
	Name                string             `bson:"name,omitempty"`
	Description         string             `bson:"description,omitempty"`
	OrganizationMembers []Member           `bson:"organization_members,omitempty"`
}

type Member struct {
	Name        string `bson:"name,omitempty"`
	Email       string `bson:"email,omitempty"`
	AccessLevel string `bson:"access_level,omitempty"`
}
