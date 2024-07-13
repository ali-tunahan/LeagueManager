package services

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"errors"
	"math/rand"
)

type LeagueService interface {
	CreateLeague(league *models.League) error
	GetLeagueByID(id uint) (*models.League, error)
	UpdateLeague(league *models.League) error
	DeleteLeague(id uint) error
	GetAllLeagues() ([]*models.League, error)
	AddTeamToLeague(leagueID, teamID uint) error
	RemoveTeamFromLeague(leagueID, teamID uint) error
	AdvanceWeek(leagueID uint) error
	ViewMatchResults(leagueID uint) ([]models.Match, error)
	EditMatchResults(leagueID, matchID uint, updatedMatch models.Match) error
	PredictChampion(leagueID uint) (models.Team, error)
	PlayAllMatches(leagueID uint) error
}

type LeagueServiceImpl struct {
	leagueRepo   repositories.LeagueRepository
	teamRepo     repositories.TeamRepository
	matchRepo    repositories.MatchRepository
	standingRepo repositories.StandingRepository
}

func NewLeagueService(leagueRepo repositories.LeagueRepository, teamRepo repositories.TeamRepository, matchRepo repositories.MatchRepository, standingRepo repositories.StandingRepository) LeagueService {
	return &LeagueServiceImpl{
		leagueRepo:   leagueRepo,
		teamRepo:     teamRepo,
		matchRepo:    matchRepo,
		standingRepo: standingRepo,
	}
}

func (s *LeagueServiceImpl) CreateLeague(league *models.League) error {
	return s.leagueRepo.CreateLeague(league)
}

func (s *LeagueServiceImpl) GetLeagueByID(id uint) (*models.League, error) {
	return s.leagueRepo.GetLeagueByID(id)
}

func (s *LeagueServiceImpl) UpdateLeague(league *models.League) error {
	return s.leagueRepo.UpdateLeague(league)
}

func (s *LeagueServiceImpl) DeleteLeague(id uint) error {
	return s.leagueRepo.DeleteLeague(id)
}

func (s *LeagueServiceImpl) GetAllLeagues() ([]*models.League, error) {
	return s.leagueRepo.GetAllLeagues()
}

func (s *LeagueServiceImpl) GetLeaguesByTeamID(teamID uint) ([]*models.League, error) {
	return s.leagueRepo.GetLeaguesByTeamID(teamID)
}

func (s *LeagueServiceImpl) AddTeamToLeague(leagueID, teamID uint) error {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return err
	}

	if len(league.Teams) >= 4 {
		return errors.New("cannot add more than 4 teams to a league")
	}

	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return err
	}

	league.Teams = append(league.Teams, *team)
	return s.leagueRepo.UpdateLeague(league)
}

func (s *LeagueServiceImpl) RemoveTeamFromLeague(leagueID, teamID uint) error {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return err
	}

	for i, team := range league.Teams {
		if team.ID == teamID {
			league.Teams = append(league.Teams[:i], league.Teams[i+1:]...)
			break
		}
	}

	return s.leagueRepo.UpdateLeague(league)
}

// AdvanceWeek advances the league to the next week and plays the matches for that week
func (s *LeagueServiceImpl) AdvanceWeek(leagueID uint) error {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return err
	}

	if len(league.Teams) != 4 {
		return errors.New("league must have exactly 4 teams to advance")
	}

	// Play matches for the current week
	matches, err := s.playMatches(league)
	if err != nil {
		return err
	}

	// Save match results
	for _, match := range matches {
		err := s.saveMatchResult(&match)
		if err != nil {
			return err
		}
	}

	// Advance the league week
	league.CurrentWeek++
	if league.CurrentWeek > 38 { // TODO write a function inside league entity instead
		return errors.New("league has already ended")
	}

	return s.leagueRepo.UpdateLeague(league)
}

// playMatches simulates the matches for the current week
func (s *LeagueServiceImpl) playMatches(league *models.League) ([]models.Match, error) {
	var matches []models.Match
	teams := league.Teams

	if len(teams) != 4 {
		return nil, errors.New("league must have exactly 4 teams to play matches")
	}

	// Example fixtures for 4 teams:
	// Week 1: A vs B, C vs D
	// Week 2: A vs C, B vs D
	// Week 3: A vs D, B vs C
	// Repeat for each set of matches

	weekFixtures := [][][2]int{
		{{0, 1}, {2, 3}},
		{{0, 2}, {1, 3}},
		{{0, 3}, {1, 2}},
	}

	weekIndex := (league.CurrentWeek - 1) % 3
	fixtures := weekFixtures[weekIndex]

	for _, fixture := range fixtures {
		homeTeam := teams[fixture[0]]
		awayTeam := teams[fixture[1]]
		homeScore, awayScore := s.simulateMatch(homeTeam, awayTeam)

		match := models.Match{
			LeagueID:      league.ID,
			HomeTeamID:    homeTeam.ID,
			AwayTeamID:    awayTeam.ID,
			HomeTeamScore: homeScore,
			AwayTeamScore: awayScore,
			Week:          league.CurrentWeek,
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (s *LeagueServiceImpl) ViewMatchResults(leagueID uint) ([]models.Match, error) {
	//TODO implement me
	panic("implement me")
}

func (s *LeagueServiceImpl) EditMatchResults(leagueID, matchID uint, updatedMatch models.Match) error {
	//TODO implement me
	panic("implement me")
}

func (s *LeagueServiceImpl) PredictChampion(leagueID uint) (models.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (s *LeagueServiceImpl) PlayAllMatches(leagueID uint) error {
	//TODO implement me
	panic("implement me")
}

// Below are helper methods for simulating matches and updating standings

// simulateMatch simulates the result of a match based on teams' strengths
func (s *LeagueServiceImpl) simulateMatch(homeTeam, awayTeam models.Team) (int, int) {

	homeAttack := homeTeam.AttackStrength
	awayDefense := awayTeam.DefenseStrength
	awayAttack := awayTeam.AttackStrength
	homeDefense := homeTeam.DefenseStrength

	homeScore := s.calculateScore(homeAttack, awayDefense)
	awayScore := s.calculateScore(awayAttack, homeDefense)

	return homeScore, awayScore
}

// calculateScore calculates the score for a team based on its attack strength and the opponent's defense strength
func (s *LeagueServiceImpl) calculateScore(attack, defense int) int {
	baseScore := rand.Intn(3) // Random base score between 0 and 2
	attackFactor := rand.Float64() * float64(attack) / 100
	defenseFactor := rand.Float64() * float64(defense) / 100

	score := baseScore + int(attackFactor*10) - int(defenseFactor*5)
	if score < 0 {
		score = 0
	}

	return score
}

// saveMatchResult saves the match result and updates the standings
func (s *LeagueServiceImpl) saveMatchResult(match *models.Match) error {
	if err := s.matchRepo.CreateMatch(match); err != nil {
		return err
	}

	// Update standings for home team
	if err := s.updateStandings(match.LeagueID, match.HomeTeamID, match.HomeTeamScore, match.AwayTeamScore); err != nil {
		return err
	}

	// Update standings for away team
	if err := s.updateStandings(match.LeagueID, match.AwayTeamID, match.AwayTeamScore, match.HomeTeamScore); err != nil {
		return err
	}

	return nil
}

// updateStandings updates the standings based on match results
func (s *LeagueServiceImpl) updateStandings(leagueID, teamID uint, teamScore, opponentScore int) error {
	standings, err := s.standingRepo.GetStandingByTeam(leagueID, teamID)
	if err != nil {
		// Create new standings if not exists
		standings = &models.Standing{
			LeagueID:       leagueID,
			TeamID:         teamID,
			Points:         0,
			Played:         0,
			Wins:           0,
			Draws:          0,
			Losses:         0,
			GoalDifference: 0,
		}
	}

	standings.GoalDifference += teamScore - opponentScore

	if teamScore > opponentScore {
		standings.Wins++
		standings.Points += 3
	} else if teamScore == opponentScore {
		standings.Draws++
		standings.Points++
	} else {
		standings.Losses++
	}

	// Standing is newly created if err is not nil
	if err != nil {
		return s.standingRepo.CreateStanding(standings)
	}
	return s.standingRepo.UpdateStanding(standings)

}
