package routes

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
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
}

func HandleCreateReaction(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := context.Background()

		// code

		return c.JSON(http.StatusCreated, nil)
	}
}
