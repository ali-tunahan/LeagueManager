package models

import "gorm.io/gorm"

// League represents a football league with teams and its current state
type League struct {
	gorm.Model
	Name        string `json:"name"`
	CurrentWeek int    `json:"current_week"`
	Teams       []Team `json:"teams" gorm:"many2many:league_teams;"`
}
