package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/ts22082/dat_board_be/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var failedToGetWidgetsMessage = map[string]interface{}{
	"message": "failed to get widgets",
}

func GetAllWidgets(c *fiber.Ctx) error {
	mongoDB := c.Locals("mongoDB").(*mongo.Database)

	if mongoDB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetWidgetsMessage)
	}

	var widgets []models.Widget

	widgetCollection := mongoDB.Collection("Widgets")

	cursor, err := widgetCollection.Find(context.Background(), bson.M{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetWidgetsMessage)
	}

	if err = cursor.All(context.Background(), &widgets); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(failedToGetWidgetsMessage)
	}

	return c.JSON(widgets)
}
