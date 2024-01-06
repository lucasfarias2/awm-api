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
	//
	//gcpServiceAccount := map[string]string{
	//	"type":                        "service_account",
	//	"project_id":                  "agree-with-me",
	//	"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
	//	"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
	//	"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
	//	"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
	//	"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-ra717%40agree-with-me.iam.gserviceaccount.com",
	//	"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
	//	"token_uri":                   "https://oauth2.googleapis.com/token",
	//	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	//	"universe_domain":             "googleapis.com",
	//}

	//gcpServiceAccountJson, err := json.Marshal(gcpServiceAccount)
	//if err != nil {
	//	log.Fatalf("Error marshalling fb config json: %s", err)
	//}

	//credentials, _ := google.CredentialsFromJSON(context.Background(), gcpServiceAccountJson, []string{"https://www.googleapis.com/auth/cloud-platform"}...)

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

	e.POST("/statement", routes.HandleCreateStatement(firestoreClient))
	e.GET("/statement", routes.HandleGetRandomStatement(firestoreClient))
	e.POST("/reaction", routes.HandleCreateReaction(firestoreClient))
	e.GET("/profile", routes.HandleGetProfileInformation(firestoreClient))
	e.GET("/stats", routes.HandleGetStats(firestoreClient))
	e.POST("/api/v1/auth/login", routes.HandleLogin(authClient))
	e.GET("/api/v1/auth/user", routes.HandleGetCurrentUser(authClient))

	e.Logger.Fatal(e.Start(":8080"))
}
