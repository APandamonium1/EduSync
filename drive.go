package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	ServiceAccountCredentials = "edusync-426009-343696fa49b1.json"
	FolderID                  = "1xN6TgYK_ldQL86fj89_vLS3A7L7BIZlj"
)

func DriveHandler(router *mux.Router) {
	router.HandleFunc("/instructor/drive", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/instructor/drive.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/api/files", func(res http.ResponseWriter, req *http.Request) {
		folderID := req.URL.Query().Get("folder_id")
		srv, err := createDriveService()
		if err != nil {
			http.Error(res, "Unable to create Drive service", http.StatusInternalServerError)
			return
		}

		files, err := listFiles(srv, folderID)
		if err != nil {
			http.Error(res, "Unable to retrieve files", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(files)
	}).Methods("GET")

	router.HandleFunc("/api/upload", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		folderID := req.FormValue("folderId")
		if folderID == "" {
			http.Error(res, "folderId is required", http.StatusBadRequest)
			return
		}

		fileName := req.FormValue("fileName")
		if fileName == "" {
			http.Error(res, "fileName is required", http.StatusBadRequest)
			return
		}

		file, _, err := req.FormFile("file")
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		srv, err := createDriveService()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		uploadedFile, err := uploadFileToDrive(srv, folderID, file, fileName)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(uploadedFile)
	}).Methods("POST")
}

func createDriveService() (*drive.Service, error) {
	ctx := context.Background()

	b, err := os.ReadFile(ServiceAccountCredentials)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	return srv, nil
}

func listFiles(service *drive.Service, folderID string) ([]*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents", folderID)
	fileList, err := service.Files.List().Q(query).PageSize(10).Fields("nextPageToken, files(id, name, mimeType)").Do()
	if err != nil {
		return nil, err
	}
	return fileList.Files, nil
}

func uploadFileToDrive(srv *drive.Service, folderID string, file io.Reader, fileName string) (*drive.File, error) {
	fileMetadata := &drive.File{
		Name:    fileName,
		Parents: []string{folderID},
	}

	// Create the file in Google Drive
	res, err := srv.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot create file: %v", err)
	}

	return res, nil
}
