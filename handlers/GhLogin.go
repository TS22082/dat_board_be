package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ts22082/dat_board_be/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var errorResponse = map[string]interface{}{
	"message": "failed",
}

func GhLogin(c *fiber.Ctx) error {
	code := c.Query("code")

	if code == "" || code == "null" {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	url := "https://github.com/login/oauth/access_token"
	ghAuthPayload := map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          code,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	ghAuthParams := utils.HTTPRequestParams{
		URL:     url,
		Method:  "POST",
		Headers: headers,
		Body:    ghAuthPayload,
	}

	ghAuthResults, _, err := utils.MakeHTTPRequest(ghAuthParams)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": "issue making request",
		})
	}

	var accessToken string

	if at, ok := ghAuthResults["access_token"].(string); ok {
		accessToken = at
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": "invalid access token",
		})
	}

	userEmailParams := map[string]interface{}{
		"URL":    "https://api.github.com/user/emails",
		"Method": "GET",
		"Headers": map[string]string{
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("token %v", accessToken),
		},
	}

	req, err := http.NewRequest(userEmailParams["Method"].(string), userEmailParams["URL"].(string), nil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	for key, value := range userEmailParams["Headers"].(map[string]string) {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var emailData []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &emailData)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"error": "cannot unmarshall this data",
		})
	}

	primaryEmail := ""
	for _, email := range emailData {
		if email["primary"].(bool) {
			primaryEmail = email["email"].(string)
			break
		}
	}

	mongoDB := c.Locals("mongoDB").(*mongo.Database)

	userCollection := mongoDB.Collection("Users")
	userFound := userCollection.FindOne(context.Background(), bson.D{{Key: "email", Value: primaryEmail}})

	type UserModel struct {
		Email string `json:"email" bson:"email"`
		ID    string `json:"_id" bson:"_id"`
	}

	var user UserModel
	err = userFound.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No user found, create a new user
			newUser := UserModel{Email: primaryEmail}
			insertResult, err := userCollection.InsertOne(context.Background(), newUser)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
					"error": "Failed to create user",
				})
			}

			newUser.ID = insertResult.InsertedID.(string)
			// Generate JWT token for the new user
			token, err := utils.GenerateJWT(newUser.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
					"error": "Failed to generate token",
				})
			}

			return c.JSON(map[string]interface{}{
				"user":  newUser,
				"token": token,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": "Error looking up user",
		})
	}

	// User found, generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": "Failed to generate token",
		})
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	return c.JSON(response)
}
