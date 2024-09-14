package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatorId string             `json:"creatorId" bson:"creatorId"`
	Title     string             `json:"title" bson:"title"`
	IsPublic  bool               `json:"isPublic" bson:"isPublic"`
	ParentId  string             `json:"parentId" bson:"parentId"`
}

func CreateItem(c *fiber.Ctx) error {
	mongoClient := c.Locals("mongoClient").(*mongo.Client)

	if mongoClient == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"message": "mongo client is nil",
		})
	}

	var itemCollection = mongoClient.Database("dat_board").Collection("Items")

	var item Item
	item.CreatorId = c.Locals("userId").(string)

	err := c.BodyParser(&item)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"message": "failed",
		})
	}

	_, err = itemCollection.InsertOne(context.Background(), item)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"message": "failed to create a new item"})
	}

	return c.JSON(map[string]interface{}{
		"message": "created",
	})
}
