package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func main() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}

	if clientID == "" || clientSecret == "" {
		fmt.Println("Error: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be set.")
		os.Exit(1)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{drive.DriveFileScope}, // Only access to files created by the app
		Endpoint:     google.Endpoint,
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n\n", authURL)

	// We cannot read stdin in non-interactive mode easily.
	// Instead, we will ask the user to run a second command with the code.
	if len(os.Args) > 1 {
		authCode := os.Args[1]
		tok, err := config.Exchange(context.Background(), authCode)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
		}

		fmt.Printf("\n--- Save this Refresh Token ---\n")
		fmt.Printf("RefreshToken: %s\n", tok.RefreshToken)
	} else {
		fmt.Println("After visiting the URL and getting the code, run this tool again with the code as an argument:")
		fmt.Println("go run tools/get_refresh_token.go <YOUR_CODE>")
	}
}
