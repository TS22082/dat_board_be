package handlers

import (
	"context"
	"fmt"
	"github.com/ts22082/dat_board_be/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	itemId := c.Params("id")

	if mongoDB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToDeleteMessage)
	}

	var itemIdHex, err = primitive.ObjectIDFromHex(itemId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToDeleteMessage)
	}

	itemCollection := mongoDB.Collection("Items")
	filter := bson.M{"_id": itemIdHex}

	_, err = itemCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToDeleteMessage)
	}

	// Delete all its descendants in a separate goroutine
	go func() {
		err := deleteSubChildren(mongoDB, itemIdHex)
		if err != nil {
			fmt.Println("Error deleting item and its descendants:", err)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "item deleted successfully",
	})
}

func deleteSubChildren(mongoDB *mongo.Database, id primitive.ObjectID) error {
	itemCollection := mongoDB.Collection("Items")

	// Find direct children of the current item
	cursor, err := itemCollection.Find(context.Background(), bson.M{"parentId": id})
	if err != nil {
		fmt.Println("Error finding children:", err)
		return err
	}

	defer cursor.Close(context.Background())

	// Iterate over each child
	for cursor.Next(context.Background()) {
		var child models.Item
		err = cursor.Decode(&child)
		if err != nil {
			fmt.Println("Error decoding child item:", err)
			return err
		}

		// Recursively delete each child and its descendants
		err = deleteSubChildren(mongoDB, child.Id)
		if err != nil {
			fmt.Println("Error deleting child and its descendants:", err)
			return err
		}

		// Finally, delete the child item itself
		_, err = itemCollection.DeleteOne(context.Background(), bson.M{"_id": child.Id})
		if err != nil {
			fmt.Println("Error deleting child item itself:", err)
			return err
		}
	}

	if err = cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return err
	}

	return nil
}
