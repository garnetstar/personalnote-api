package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"personalnote.eu/simple-go-api/utils"
)

// UploadHandler handles file uploads to Google Drive
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	userID, authenticated := checkAuth(w, r)
	if !authenticated {
		return
	}

	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is allowed")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request", "Could not parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid file", "No file provided")
		return
	}
	defer file.Close()

	// Debug logs
	log.Println("Attempting to initialize Drive service...")

	// Get Service Account credentials
	// Option 1: From environment variable (JSON content)
	credsJSON := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JSON")
	// Option 2: From file path
	credsFile := os.Getenv("GOOGLE_SERVICE_ACCOUNT_FILE")

	// Option 3: From Refresh Token (for personal accounts)
	refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	log.Printf("Env var GOOGLE_SERVICE_ACCOUNT_FILE: '%s'", credsFile)

	var driveService *drive.Service
	var serviceErr error

	ctx := context.Background()

	if refreshToken != "" && clientID != "" && clientSecret != "" {
		log.Println("Using Refresh Token for OAuth 2.0")
		config := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{drive.DriveFileScope},
		}
		token := &oauth2.Token{RefreshToken: refreshToken}
		tokenSource := config.TokenSource(ctx, token)
		driveService, serviceErr = drive.NewService(ctx, option.WithTokenSource(tokenSource))
	} else if credsJSON != "" {
		log.Println("Using credsJSON")
		driveService, serviceErr = drive.NewService(ctx, option.WithCredentialsJSON([]byte(credsJSON)), option.WithScopes(drive.DriveScope))
	} else if credsFile != "" {
		log.Printf("Using credsFile: %s", credsFile)
		// Check if file exists
		if _, err := os.Stat(credsFile); os.IsNotExist(err) {
			log.Printf("❌ Credentials file does not exist at path: %s", credsFile)
			serviceErr = fmt.Errorf("credentials file not found: %s", credsFile)
		} else {
			driveService, serviceErr = drive.NewService(ctx, option.WithCredentialsFile(credsFile), option.WithScopes(drive.DriveScope))
		}
	} else {
		log.Println("❌ Google Service Account credentials not found")
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Configuration error", "Google Drive integration is not configured")
		return
	}

	if serviceErr != nil {
		log.Printf("❌ Failed to create Drive service: %v", serviceErr)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Service error", "Failed to connect to Google Drive")
		return
	}

	// Create file metadata
	driveFile := &drive.File{
		Name: header.Filename,
	}

	folderID := os.Getenv("GOOGLE_DRIVE_FOLDER_ID")
	if folderID != "" {
		driveFile.Parents = []string{folderID}
	} else {
		log.Println("⚠️ GOOGLE_DRIVE_FOLDER_ID not set. Attempting to upload to root (likely to fail for Service Accounts).")
	}

	// Upload file
	log.Printf("📤 Uploading file '%s' (Size: %d bytes) for user %d", header.Filename, header.Size, userID)

	uploadedFile, err := driveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		log.Printf("❌ Failed to upload file to Drive: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Upload error", "Failed to upload file to Google Drive")
		return
	}

	log.Printf("✅ File uploaded successfully. ID: %s", uploadedFile.Id)

	// Return success response
	response := map[string]interface{}{
		"message":     "File uploaded successfully",
		"fileId":      uploadedFile.Id,
		"name":        uploadedFile.Name,
		"mimeType":    uploadedFile.MimeType,
		"webViewLink": uploadedFile.WebViewLink,
	}

	utils.SendJSONResponse(w, http.StatusCreated, response)
}
