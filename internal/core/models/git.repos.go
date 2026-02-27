package models

import (
	"time"
)
type Repo struct {
	Name 		string 	`json:"name"`
	Language 	string 	`json:"language"`
}
type GitHubUser struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GitHubLanguageReport struct {
	Username  string         `json:"username"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Repos     []Repo         `json:"repos"`
	Languages map[string]int `json:"languages"`
	UpdatedAt time.Time      `json:"updated_at" bson:"updated_at"`
}

type GitHubUserSearchResponse struct {
	Items []struct {
		Login string `json:"login"`
	} `json:"items"`
}