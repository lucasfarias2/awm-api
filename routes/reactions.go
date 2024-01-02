package routes

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

type ReactionType string

const (
	Agreed    ReactionType = "agreed"
	Disagreed ReactionType = "disagreed"
	Skipped   ReactionType = "skipped"
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
	UserID      string       `json:"user_id"`
	Reaction    ReactionType `json:"reaction"`
	CreatedAt   time.Time    `json:"created_at"`
}

func HandleCreateReaction(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		location, err := time.LoadLocation("CET")

		var req Reaction
		if err := c.Bind(&req); err != nil {
			log.Fatalf("Failed to bind request: %v", err)
		}

		newReac, _, err := client.Collection("reactions").Add(ctx, Reaction{
			StatementID: req.StatementID,
			UserID:      req.UserID,
			Reaction:    req.Reaction,
			CreatedAt:   time.Now().In(location),
		})
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return c.JSON(http.StatusCreated, newReac)
	}
}
