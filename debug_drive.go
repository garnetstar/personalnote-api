package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	credsFile := "keys/quickstart-1549817042430-d5f603eed637.json"
	folderID := "1_Ya2XCpaZNf5VlKLO9kZycnqX41htOdC"

	fmt.Printf("Checking folder %s using creds %s\n", folderID, credsFile)

	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credsFile), option.WithScopes(drive.DriveScope))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// 1. Check About (Quota info)
	about, err := srv.About.Get().Fields("storageQuota,user").Do()
	if err != nil {
		log.Printf("Error getting about info: %v", err)
	} else {
		fmt.Printf("Service Account User: %s\n", about.User.EmailAddress)
		fmt.Printf("Quota Limit: %d\n", about.StorageQuota.Limit)
		fmt.Printf("Quota Usage: %d\n", about.StorageQuota.Usage)
	}

	// 2. Check Folder Access
	f, err := srv.Files.Get(folderID).Fields("id, name, capabilities, owners, shared").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve file: %v", err)
	}

	fmt.Printf("Folder Name: %s\n", f.Name)
	fmt.Printf("Can Add Children: %v\n", f.Capabilities.CanAddChildren)
	fmt.Printf("Is Shared: %v\n", f.Shared)

	owners, _ := json.Marshal(f.Owners)
	fmt.Printf("Owners: %s\n", string(owners))
}
