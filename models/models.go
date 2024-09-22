package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	CreatorId string `json:"creatorId" bson:"creatorId"`
	Id        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title"`
	IsPublic  bool   `json:"isPublic" bson:"isPublic"`
	ParentId  string `json:"parentId" bson:"parentId"`
}

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Email string             `bson:"email,omitempty"`
}
