package models

// Standings represents the standings of a team in a specific league
type Standings struct {
	ID             uint `json:"id" gorm:"primaryKey"`
	LeagueID       uint `json:"league_id"`
	TeamID         uint `json:"team_id"`
	Points         int  `json:"points"`
	Played         int  `json:"played"`
	Wins           int  `json:"wins"`
	Draws          int  `json:"draws"`
	Losses         int  `json:"losses"`
	GoalDifference int  `json:"goal_difference"`
}
