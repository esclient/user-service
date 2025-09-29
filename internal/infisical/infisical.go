package infisical

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	InfisicalURL = "https://us.infisical.com"
	SecretPathUserService = "/user-service"
) 

type RetrievedSecret struct {
	Key   string `json:"secretKey"`
	Value string `json:"secretValue"`
}

type Client struct {
	BaseURL   string
	AuthToken string
	http      *http.Client
}

func NewClient(baseURL string, authToken string) *Client {
	return &Client{
		BaseURL:   baseURL,
		AuthToken: authToken,
		http:      &http.Client{},
	}
}

func (c *Client) GetSecretsV4(ctx context.Context, projectId string, environment string, secretsPath string) ([]RetrievedSecret, error) {
	url := fmt.Sprintf("%s/api/v4/secrets?projectId=%s&environment=%s&secretPath=%s",
		c.BaseURL, projectId, environment, secretsPath)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("infisical v4 error: %s", string(body))
	}

	var data struct {
		Secrets []RetrievedSecret `json:"secrets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Secrets, nil
}

// Функция взятия конкретного ключа, а не списка
func (c *Client) GetSecret(ctx context.Context, projectId string, environment string, secretsPath string, key string) (*RetrievedSecret, error) {
	secrets, err := c.GetSecretsV4(ctx, projectId, environment, secretsPath)
	if err != nil {
		return nil, err
	}
	for _, s := range secrets {
		if s.Key == key {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("secret %s not found", key)
}
