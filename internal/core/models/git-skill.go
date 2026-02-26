package models

import "time"

type UserSkills struct {
	Username string 	`bson:"username" json:"username"`
	Skills   map[string]float64 `bson:"skills" json:"skills"`
	LastUpdated time.Time 	`bson:"last_updated" json:"last_updated"`
}