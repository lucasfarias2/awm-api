package routes

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

// Statement represents the structure of our resource
type Statement struct {
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	CreatedAt time.Time
}

func HandleCreateStatement(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		location, err := time.LoadLocation("CET")

		var req Statement
		if err := c.Bind(&req); err != nil {
			log.Fatalf("Failed to bind request: %v", err)
		}

		newS, _, err := client.Collection("statements").Add(ctx, Statement{
			Text:      req.Text,
			UserID:    req.UserID,
			CreatedAt: time.Now().In(location),
		})
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return c.JSON(http.StatusCreated, interface{}(newS))
	}
}

func HandleGetRandomStatement(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusOK, nil)
	}
}
