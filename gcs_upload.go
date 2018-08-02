package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
		"golang.org/x/net/context"

	"cloud.google.com/go/storage"

	"google.golang.org/api/option"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)


type bucketConfig struct {
	keyPath string
	name string
}

type artefact struct {
	folderName string
	uploadFilePath string
	uploadedFileName string
}

func main() {
	bucketConfig := bucketConfig{
		os.Getenv("GCS_SERVICE_ACCOUNT_JSON_KEY_URL"),
		os.Getenv("BUCKET_NAME"),
	}

	artefact := artefact {
		os.Getenv("FOLDER_NAME"),
		os.Getenv("UPLOAD_FILE_PATH"),
		os.Getenv("UPLOADED_FILE_NAME"),
	}

	log.Debugf("KeyPath => %s", bucketConfig.keyPath)
	log.Debugf("bucketName => %s", bucketConfig.name)
	log.Debugf("folderName => %s", artefact.folderName)
	log.Debugf("uploadFilePath => %s", artefact.uploadFilePath)
	log.Debugf("uploadedFileName => %s", artefact.uploadedFileName)


	localKeyPath := downloadKeyFile(bucketConfig.keyPath)
	log.Debugf("local localKeyPath => %s", localKeyPath)
	setGoogleCredentials(localKeyPath)

	context := context.Background()
	client := createClient(context, localKeyPath)

	uploadFile(context, client, artefact, bucketConfig)
	closeClient(client)
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func downloadFile(downloadURL, targetPath string) error {
	outFile, err := os.Create(targetPath)
	if err != nil {
		failf("Failed to create (%s), error: %s", targetPath, err)
	}
	defer func() {
		if err = outFile.Close(); err != nil {
			log.Warnf("Failed to close (%s), error: %s", targetPath, err)
		}
	}()

	resp, err := http.Get(downloadURL)
	if err != nil {
		failf("Failed to download from (%s), error: %s", downloadURL, err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Warnf("Failed to close (%s) body", downloadURL)
		}
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		failf("Failed to download from (%s), error: %s", downloadURL, err)
	}

	return nil
}

func createClient(context context.Context, keyPath string) *storage.Client {
	client, err := storage.NewClient(context, option.WithCredentialsFile(keyPath))

	if err != nil {
		failf("Failed to create new storage client, error: %s", err)
	}

	return client;
}

func closeClient(client *storage.Client) {
	if err := client.Close(); err != nil {
		failf("Failed to close storage client, error: %s", err)
	}
}

func uploadFile(context context.Context, client *storage.Client, artefact artefact, bucketConfig bucketConfig) {
	file, err := os.Open(artefact.uploadFilePath)

	if err != nil {
		failf("File (%s) does not exist, error: %s", artefact.uploadFilePath, err)
	}

	defer file.Close()

	uploadPath := filePath(artefact.folderName, artefact.uploadedFileName)
	bkt := client.Bucket(bucketConfig.name)
	writer := bkt.Object(uploadPath).NewWriter(context)

	copyFileToWriter(writer, file, uploadPath)
	closeWriter(writer)
}

func closeWriter(writer *storage.Writer) {
	if err := writer.Close(); err != nil {
		failf("Failed to close writer, error: %s", err)
	}
}

func copyFileToWriter(writer *storage.Writer, file *os.File, uploadFilePath string) {
	if _, err := io.Copy(writer, file); err != nil {
		failf("File (%s) does not exist, error: %s", uploadFilePath, err)
	}
}

func filePath(folderName string, uploadedFileName string) string {
	if folderName != "" {
		uploadedFileName = folderName + "/" + uploadedFileName
	}
	return uploadedFileName
}

func downloadKeyFile(keyPath string) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__google-cloud-storage__")

	if err != nil {
		failf("Failed to create tmp dir, error: %s", err)
	}

	targetPath := filepath.Join(tmpDir, "key.json")

	if err := downloadFile(keyPath, targetPath); err != nil {
		failf("Failed to download json key file, error: %s", err)
	}

	return targetPath
}

func setGoogleCredentials(keyPath string) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", keyPath)
}
