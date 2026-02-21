package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

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

	// Get Service Account credentials
	// Option 1: From environment variable (JSON content)
	credsJSON := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JSON")
	// Option 2: From file path
	credsFile := os.Getenv("GOOGLE_SERVICE_ACCOUNT_FILE")

	var driveService *drive.Service
	var serviceErr error

	ctx := context.Background()

	if credsJSON != "" {
		driveService, serviceErr = drive.NewService(ctx, option.WithCredentialsJSON([]byte(credsJSON)), option.WithScopes(drive.DriveScope))
	} else if credsFile != "" {
		driveService, serviceErr = drive.NewService(ctx, option.WithCredentialsFile(credsFile), option.WithScopes(drive.DriveScope))
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
