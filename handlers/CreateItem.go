package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	CreatorId string `json:"creatorId" bson:"creatorId"`
	Id        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title"`
	IsPublic  bool   `json:"isPublic" bson:"isPublic"`
	ParentId  string `json:"parentId" bson:"parentId"`
}

var failedToCreateMessage = map[string]interface{}{
	"message": "failed to create a new item",
}

func CreateItem(c *fiber.Ctx) error {
	mongoDB := c.Locals("mongoDB").(*mongo.Database)

	if mongoDB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	var itemCollection = mongoDB.Collection("Items")

	var item Item
	item.CreatorId = c.Locals("userId").(string)

	err := c.BodyParser(&item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	res, err := itemCollection.InsertOne(context.Background(), item)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	item.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return c.JSON(map[string]interface{}{
		"message": item,
	})
}
