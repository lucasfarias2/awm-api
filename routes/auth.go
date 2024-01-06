package routes

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func HandleGetCurrentUser(auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()

		token, err := auth.VerifySessionCookieAndCheckRevoked(ctx, c.Request().Header.Get("session"))
		if err != nil {
			log.Printf("Failed to authenticate request: %v", err)
			return err
		}

		user, err := auth.GetUser(ctx, token.UID)
		if err != nil {
			log.Printf("Failed to authenticate request: %v", err)
			return err
		}

		// Create an instance of UserResponse and populate it
		response := map[string]string{
			"id":    user.UID,
			"email": user.Email,
		}

		return c.JSON(http.StatusOK, response)
	}
}

type LoginRequest struct {
	Token string `json:"token"`
}

func HandleLogin(auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()

		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("Failed to bind request: %v", err)
			return err
		}

		expiresIn := time.Hour * 24 * 14

		token, err := auth.SessionCookie(ctx, req.Token, expiresIn)
		if err != nil {
			log.Printf("Failed to create a session cookie: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, token)
	}
}
