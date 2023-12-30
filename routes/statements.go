package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Statement represents the structure of our resource
type Statement struct {
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

func HandleCreateStatement(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusCreated, nil)
	}
}

func HandleGetRandomStatement(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusOK, nil)
	}
}
