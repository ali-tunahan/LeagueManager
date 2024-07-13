package dto

// TeamPrediction represents the predicted probability of a team winning the league
type TeamPrediction struct {
	LeagueID       uint    `json:"league_id"`
	TeamID         uint    `json:"team_id"`
	TeamName       string  `json:"team_name"`
	WinProbability float64 `json:"win_probability"`
}
