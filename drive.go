package main

import (
	"context"
	"encoding/json"
	"fmt"
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
)

func DriveHandler(router *mux.Router) {
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

	router.HandleFunc("/api/media-upload", func(res http.ResponseWriter, req *http.Request) {
		// Initialize Google Drive client
		srv, err := createDriveService()
		if err != nil {
			http.Error(res, "Failed to initialize Drive client", http.StatusInternalServerError)
			return
		}

		err = req.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(res, "Failed to parse form", http.StatusBadRequest)
			return
		}

		files := req.MultipartForm.File["files"]
		folderIds := req.MultipartForm.Value["folderIds"]
		studentNames := req.MultipartForm.Value["studentNames"]

		if len(files) != len(folderIds) || len(folderIds) != len(studentNames) {
			http.Error(res, "Mismatch between files, folder IDs, and student names", http.StatusBadRequest)
			return
		}

		var results []string

		for i, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(res, fmt.Sprintf("Failed to open file: %v", err), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			folderID := folderIds[i]
			studentName := studentNames[i]

			_, err = uploadFileToDrive(srv, folderID, file, fileHeader.Filename)
			if err != nil {
				results = append(results, fmt.Sprintf("Failed to upload %s for %s: %v", fileHeader.Filename, studentName, err))
			} else {
				results = append(results, fmt.Sprintf("Successfully uploaded %s for %s", fileHeader.Filename, studentName))
			}
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(results)
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
