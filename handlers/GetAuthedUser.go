package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ts22082/dat_board_be/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAuthedUser gets the authenticated user from the database
// and returns it as a JSON response
// Parameters:
// - c: The Fiber context, used to retrieve the MongoDB database and request information.
//
// Returns:
// - An error if the item creation fails, or a JSON response with the created item.
func GetAuthedUser(c *fiber.Ctx) error {

	userId := c.Locals("userId").(string)
	mongoDB := c.Locals("mongoDB").(*mongo.Database)

	objID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	userCollection := mongoDB.Collection("Users")
	err = userCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: objID}}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(map[string]interface{}{
				"message": "There is no documents that match",
			})
		} else {
			log.Fatal(err)
		}
		return c.JSON(map[string]interface{}{
			"message": "There was a problem finding this document",
		})
	}

	log.Printf("user: %v", user)

	return c.JSON(map[string]interface{}{
		"user": user,
	})
}
