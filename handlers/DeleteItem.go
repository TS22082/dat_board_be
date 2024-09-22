package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var failedToDeleteMessage = map[string]interface{}{
	"message": "failed to delete item",
}

// DeleteItem handles the deletion of an item.
// It expects the MongoDB instance and the user ID to be passed in the request context.
// The function deletes the item from the MongoDB collection and returns a JSON response with the success message.
//
// Parameters:
// - c: The Fiber context, used to retrieve the MongoDB database and request information.
//
// Returns:
// - An error if the item deletion fails, or a JSON response with the success message.

func DeleteItem(c *fiber.Ctx) error {
	mongoDB := c.Locals("mongoDB").(*mongo.Database)
	userId := c.Locals("userId").(string)
	itemId := c.Params("id")

	if mongoDB == nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToDeleteMessage)
	}

	itemCollection := mongoDB.Collection("Items")
	filter := bson.M{"_id": itemId, "creatorId": userId}

	_, err := itemCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(failedToDeleteMessage)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "item deleted successfully",
	})
}
