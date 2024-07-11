package models

import "gorm.io/gorm"

type League struct {
	gorm.Model
	Name        string     `json:"name"`
	CurrentWeek int        `json:"current_week"`
	Teams       []Team     `json:"teams" gorm:"many2many:league_teams;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Matches     []Match    `json:"matches" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Standings   []Standing `json:"standings" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
