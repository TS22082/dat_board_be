package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ts22082/dat_board_be/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Item represents an item structure with fields corresponding to the JSON and BSON object structure.

var failedToCreateMessage = map[string]interface{}{
	"message": "failed to create a new item",
}

// CreateItem handles the creation of a new item.
// It expects the MongoDB instance and the user ID to be passed in the request context.
// The function reads the item details from the request body, inserts the new item into the MongoDB collection,
// and returns a JSON response with the created item or an error message if the operation fails.
//
// Parameters:
// - c: The Fiber context, used to retrieve the MongoDB database and request information.
//
// Returns:
// - An error if the item creation fails, or a JSON response with the created item.
func CreateItem(c *fiber.Ctx) error {
	mongoDB := c.Locals("mongoDB").(*mongo.Database)

	if mongoDB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	var itemCollection = mongoDB.Collection("Items")
	var item models.Item
	var stringFromLocals = c.Locals("userId").(string)

	var creatorId, err = primitive.ObjectIDFromHex(stringFromLocals)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	item.CreatorId = creatorId

	currentTime := time.Now().UTC().Format(time.RFC3339)

	item.CreatedAt = currentTime
	item.UpdatedAt = currentTime

	err = c.BodyParser(&item)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	res, err := itemCollection.InsertOne(context.Background(), item)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToCreateMessage)
	}

	item.Id = res.InsertedID.(primitive.ObjectID)

	return c.JSON(map[string]interface{}{
		"message": item,
	})
}
