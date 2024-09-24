package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/ts22082/dat_board_be/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetItem handles the retrieval of an item.
// It expects the MongoDB instance and the user ID to be passed in the request context.
// The function retrieves the item from the MongoDB collection and returns a JSON response with the item data.
//
// Parameters:
// - c: The Fiber context, used to retrieve the MongoDB database and request information.
//
// Returns:
// - An error if the item retrieval fails, or a JSON response with the item data.

func GetItem(c *fiber.Ctx) error {
	mongoDB := c.Locals("mongoDB").(*mongo.Database)
	id := c.Params("id")

	if mongoDB == nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	var itemCollection = mongoDB.Collection("Items")
	var item models.Item

	itemAsHex, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	err = itemCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: itemAsHex}}).Decode(&item)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	return c.Status(fiber.StatusOK).JSON(item)
}
