package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitHubClient struct {
	BaseURL string
	Token   string
}

type Repo struct {
	Name string `json:"name"`
}

func (g *GitHubClient) GetRepos(ctx context.Context, username string) ([]Repo, error) {
	url := fmt.Sprintf("%s/users/%s/repos", g.BaseURL, username)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	if g.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repo
	err = json.NewDecoder(resp.Body).Decode(&repos)
	return repos, err
}