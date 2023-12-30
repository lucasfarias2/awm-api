package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleGetProfileInformation(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusOK, nil)
	}
}
