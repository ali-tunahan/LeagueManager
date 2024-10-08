package models

import "gorm.io/gorm"

// Standing represents the standings of a team in a specific league
type Standing struct {
	gorm.Model
	LeagueID       uint `json:"league_id"`
	TeamID         uint `json:"team_id"`
	Points         int  `json:"points"`
	Played         int  `json:"played"`
	Wins           int  `json:"wins"`
	Draws          int  `json:"draws"`
	Losses         int  `json:"losses"`
	GoalDifference int  `json:"goal_difference"`
}
