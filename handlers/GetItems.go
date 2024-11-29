package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/ts22082/dat_board_be/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var failedToGetMessage = map[string]interface{}{
	"message": "failed to get items",
}

// GetItems handles the retrieval of all items created by the authenticated user.
// It retrieves the items from the MongoDB collection and returns them as a JSON response.
// Parameters:
// - c: The Fiber context, used to retrieve the MongoDB database and request information.
//
// Returns:
// - An error if the item retrieval fails, or a JSON response with the retrieved items.
func GetItems(c *fiber.Ctx) error {
	mongoDB := c.Locals("mongoDB").(*mongo.Database)
	userIdStr := c.Locals("userId").(string)
	parentIdStr := c.Query("parentId")

	if mongoDB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	var items []models.Item

	itemCollection := mongoDB.Collection("Items")
	uidAsHex, err := primitive.ObjectIDFromHex(userIdStr)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	filter := bson.M{"creatorId": uidAsHex, "parentId": nil}

	if parentIdStr != "" {
		parentIdAsHex, err := primitive.ObjectIDFromHex(parentIdStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
		}
		filter = bson.M{"creatorId": uidAsHex, "parentId": parentIdAsHex}
	}

	cursor, err := itemCollection.Find(context.Background(), filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	if err = cursor.All(context.Background(), &items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	return c.JSON(items)
}
