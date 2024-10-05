package handlers

import (
	"context"

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
	userId := c.Locals("userId").(string)
	parentId := c.Query("parentId")

	if mongoDB == nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	var items []models.Item

	itemCollection := mongoDB.Collection("Items")
	filter := bson.M{"creatorId": userId, "parentId": parentId}

	cursor, err := itemCollection.Find(context.Background(), filter)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	if err = cursor.All(context.Background(), &items); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToGetMessage)
	}

	return c.JSON(items)
}
