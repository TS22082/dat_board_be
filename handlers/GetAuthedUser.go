package handlers

import (
	"context"
	"errors"
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

	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.JSON(map[string]interface{}{
			"message": "There is no document that matches",
		})
	}

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Issue looking up document",
		})
	}

	return c.JSON(map[string]interface{}{
		"user": user,
	})
}
