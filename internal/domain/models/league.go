package models

// League represents a football league with teams and its current state
type League struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	CurrentWeek int    `json:"current_week"`
	Teams       []Team `json:"teams" gorm:"many2many:league_teams;"`
}
