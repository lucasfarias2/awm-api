package routes

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

type ReactionType string

const (
	Agreed    ReactionType = "agree"
	Disagreed ReactionType = "disagree"
	Skipped   ReactionType = "skip"
)

// Validate checks if the reaction type is valid
func (r ReactionType) Validate() error {
	switch r {
	case Agreed, Disagreed, Skipped:
		return nil
	default:
		return fmt.Errorf("invalid reaction type: %s", r)
	}
}

type Reaction struct {
	StatementID string       `json:"statement_id"`
	Reaction    ReactionType `json:"reaction"`
	CreatedAt   time.Time    `json:"created_at"`
}

type NewReactionRequest struct {
	StatementID string
	Reaction    ReactionType
	UserID      string
	CreatedAt   time.Time
}

func HandleCreateReaction(client *firestore.Client, auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		location, err := time.LoadLocation("CET")

		var req Reaction
		if err := c.Bind(&req); err != nil {
			log.Fatalf("Failed to bind request: %v", err)
		}

		session := c.Request().Header.Get("session")
		token, err := auth.VerifySessionCookieAndCheckRevoked(ctx, session)
		if err != nil {
			log.Printf("Failed to authenticate request: %v", err)
			return err
		}
		user, err := auth.GetUser(ctx, token.UID)
		if err != nil {
			log.Printf("Failed to authenticate request: %v", err)
			return err
		}

		newReac, _, err := client.Collection("reactions").Add(ctx, NewReactionRequest{
			StatementID: req.StatementID,
			UserID:      user.UID,
			Reaction:    req.Reaction,
			CreatedAt:   time.Now().In(location),
		})
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return c.JSON(http.StatusCreated, newReac)
	}
}
