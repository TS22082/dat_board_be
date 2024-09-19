package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Email string             `bson:"email,omitempty"`
}

// GetAuthedUser gets the authenticated user from the database
func GetAuthedUser(c *fiber.Ctx) error {

	userId := c.Locals("userId").(string)
	mongoDB := c.Locals("mongoDB").(*mongo.Database)

	objID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Fatal(err)
	}

	var user User
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
