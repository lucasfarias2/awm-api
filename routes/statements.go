package routes

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/iterator"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type StatementRequest struct {
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type StatementResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleCreateStatement(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		location, err := time.LoadLocation("CET")

		var req StatementRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("Failed to bind request: %v", err)
		}

		newS, _, err := client.Collection("statements").Add(ctx, StatementRequest{
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

type Response struct {
	Statement           StatementResponse `json:"statement"`
	AgreedPercentage    float64           `json:"agreed_percentage"`
	DisagreedPercentage float64           `json:"disagreed_percentage"`
}

func HandleGetRandomStatement(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()

		userId := c.QueryParams().Get("user_id")

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

			var agreedCount, disagreedCount int

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

			for _, reaction := range statsReactions {
				if reaction.StatementID == selectedStatement.ID {
					switch reaction.Reaction {
					case Agreed:
						agreedCount++
					case Disagreed:
						disagreedCount++
					}
				}
			}

			totalReactions := agreedCount + disagreedCount
			var agreedPercentage, disagreedPercentage float64
			if totalReactions > 0 {
				agreedPercentage = float64(agreedCount) / float64(totalReactions) * 100
				disagreedPercentage = float64(disagreedCount) / float64(totalReactions) * 100
			}

			response := Response{
				Statement:           selectedStatement,
				AgreedPercentage:    math.Round(agreedPercentage*100) / 100,
				DisagreedPercentage: math.Round(disagreedPercentage*100) / 100,
			}

			return c.JSON(http.StatusOK, response)
		} else {
			return c.JSON(http.StatusOK, nil)
		}
	}
}
