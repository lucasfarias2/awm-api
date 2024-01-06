package main

import (
	"awm-api/routes"
	"context"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	e := echo.New()

	config := &firebase.Config{ProjectID: "agree-with-me"}

	app, err := firebase.NewApp(context.Background(), config, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx := context.Background()

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}

	e.POST("/api/v1/statement", routes.HandleCreateStatement(firestoreClient, authClient))
	e.GET("/api/v1/statement", routes.HandleGetUserStatements(firestoreClient, authClient))
	e.GET("/api/v1/feed", routes.HandleGetRandomStatement(firestoreClient, authClient))
	e.POST("/api/v1/reaction", routes.HandleCreateReaction(firestoreClient, authClient))
	e.GET("/api/v1/profile", routes.HandleGetProfileInformation(firestoreClient, authClient))
	e.GET("/api/v1/stats", routes.HandleGetStats(firestoreClient))
	e.POST("/api/v1/auth/login", routes.HandleLogin(authClient))
	e.GET("/api/v1/auth/user", routes.HandleGetCurrentUser(authClient))

	e.Logger.Fatal(e.Start(":8080"))
}
