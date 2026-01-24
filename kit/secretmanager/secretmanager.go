package secretmanager

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type Config struct {
	ProjectID string
}
type Client struct {
	*secretmanager.Client
	projectID string
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secretmanager.retrievekid failed to create secret manager client: %w", err)
	}

	return &Client{Client: client, projectID: cfg.ProjectID}, nil
}

func (c *Client) GetSecret(ctx context.Context, name string) ([]byte, error) {
	req := secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.projectID, name),
	}

	result, err := c.AccessSecretVersion(ctx, &req)
	if err != nil {
		return []byte{}, fmt.Errorf("secretmanager.getsecret failed to access secret [%s]: %w", req.Name, err)
	}

	return result.Payload.Data, nil
}
