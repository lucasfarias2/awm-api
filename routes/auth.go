package routes

import (
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleGetCurrentUser(auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusCreated, nil)
	}
}

func HandleLogin(auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusCreated, nil)
	}
}
