package models

// Match represents a match between two teams in a specific league
type Match struct {
	ID            uint `json:"id" gorm:"primaryKey"`
	LeagueID      uint `json:"league_id"`
	HomeTeamID    uint `json:"home_team_id"`
	AwayTeamID    uint `json:"away_team_id"`
	HomeTeamScore int  `json:"home_team_score"`
	AwayTeamScore int  `json:"away_team_score"`
	Week          int  `json:"week"`
}
