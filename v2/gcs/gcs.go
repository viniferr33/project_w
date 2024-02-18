package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"project_w/v2/config"
	"project_w/v2/filehandler"
)

func UploadFile(file filehandler.File, bucket string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	f, err := os.Open(file.Filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, config.GCS_TIMEOUT_S)
	defer cancel()

	o := client.Bucket(bucket).Object(file.Name)

	o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return "", nil
	}
	if err := wc.Close(); err != nil {
		return "", nil
	}

	return fmt.Sprintf("gs://%s/%s", bucket, file.Name), nil
}
