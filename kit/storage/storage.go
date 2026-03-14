package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type Client struct {
	bucket string
	gcs    *storage.Client
}

func New(ctx context.Context, bucket string) (*Client, error) {
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.New: %w", err)
	}
	return &Client{bucket: bucket, gcs: c}, nil
}

// UploadPublic sube un objeto al bucket y lo hace públicamente accesible.
// Retorna la URL pública del objeto.
func (c *Client) UploadPublic(ctx context.Context, objectName string, data io.Reader, contentType string) (string, error) {
	obj := c.gcs.Bucket(c.bucket).Object(objectName)

	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, data); err != nil {
		wc.Close()
		return "", fmt.Errorf("storage.UploadPublic: io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("storage.UploadPublic: writer.Close: %w", err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("storage.UploadPublic: acl.Set: %w", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucket, objectName)
	return url, nil
}