package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	CreatorId primitive.ObjectID `json:"creatorId" bson:"creatorId,omitempty"`
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	IsPublic  bool               `json:"isPublic" bson:"isPublic"`
	ParentId  primitive.ObjectID `json:"parentId" bson:"parentId,omitempty"`
	CreatedAt string             `json:"createdAt" bson:"createdAt"`
	UpdatedAt string             `json:"updatedAt" bson:"updatedAt"`
}

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Email string             `bson:"email,omitempty"`
}
