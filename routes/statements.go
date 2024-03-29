package routes

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/iterator"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type StatementRequest struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStatementRequest struct {
	Text      string
	UserID    string
	CreatedAt time.Time
}

type StatementResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleGetUserStatements(client *firestore.Client, auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()

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

		userId := user.UID

		iterStatements := client.Collection("statements").Where("UserID", "==", userId).Limit(1000).Documents(ctx)
		defer iterStatements.Stop()

		// iterate over statements and push to array and return that at the end
		var statements []StatementResponse
		for {
			doc, err := iterStatements.Next()
			if errors.Is(err, iterator.Done) {
				break // Exit the loop when all documents have been read
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			var statement StatementResponse
			if err := doc.DataTo(&statement); err != nil {
				log.Fatalf("Failed to read document: %v", err)
			}

			statement.ID = doc.Ref.ID

			statements = append(statements, statement)
		}

		return c.JSON(http.StatusOK, statements)
	}
}

func HandleCreateStatement(client *firestore.Client, auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		location, err := time.LoadLocation("CET")

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

		var req StatementRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("Failed to bind request: %v", err)
		}

		newS, _, err := client.Collection("statements").Add(ctx, CreateStatementRequest{
			Text:      req.Text,
			UserID:    user.UID,
			CreatedAt: time.Now().In(location),
		})
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return c.JSON(http.StatusCreated, interface{}(newS))
	}
}

type Response struct {
	Statement StatementResponse `json:"statement"`
}

func HandleGetRandomStatement(client *firestore.Client, auth *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()

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

		userId := user.UID

		iterStatements := client.Collection("statements").Where("UserID", "!=", userId).Limit(1000).Documents(ctx)
		defer iterStatements.Stop()

		iterReactions := client.Collection("reactions").Where("UserID", "==", userId).Limit(1000).Documents(ctx)
		defer iterReactions.Stop()

		reactedStatementIDs := make(map[string]struct{})
		var reactions []Reaction
		for {
			doc, err := iterReactions.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			var reaction Reaction
			if err := doc.DataTo(&reaction); err != nil {
				log.Fatalf("Failed to read document: %v", err)
			}

			reactions = append(reactions, reaction)
			reactedStatementIDs[reaction.StatementID] = struct{}{}
		}

		var statements []StatementResponse
		for {
			doc, err := iterStatements.Next()
			if errors.Is(err, iterator.Done) {
				break // Exit the loop when all documents have been read
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			var statement StatementResponse
			if err := doc.DataTo(&statement); err != nil {
				log.Fatalf("Failed to read document: %v", err)
			}

			statement.ID = doc.Ref.ID

			if _, exists := reactedStatementIDs[statement.ID]; !exists {
				statements = append(statements, statement)
			}
		}

		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		minV := 0
		maxV := len(statements)

		if maxV > minV {
			randN := r.Intn(maxV-minV) + minV
			selectedStatement := statements[randN]

			itStatsReactions := client.Collection("reactions").Where("StatementID", "==", selectedStatement.ID).Limit(1000).Documents(ctx)
			defer itStatsReactions.Stop()

			var statsReactions []Reaction
			for {
				doc, err := itStatsReactions.Next()
				if errors.Is(err, iterator.Done) {
					break
				}
				if err != nil {
					log.Printf("Failed to iterate: %v", err)
				}

				var statReaction Reaction
				if err := doc.DataTo(&statReaction); err != nil {
					log.Printf("Failed to read document: %v", err)
				}

				statsReactions = append(statsReactions, statReaction)
			}

			response := Response{
				Statement: selectedStatement,
			}

			return c.JSON(http.StatusOK, response)
		} else {
			return c.JSON(http.StatusOK, nil)
		}
	}
}
