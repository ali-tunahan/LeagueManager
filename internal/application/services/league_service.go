package services

import (
	dto "LeagueManager/internal/domain/dtos"
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"errors"
	"math/rand"
	"sort"
)

type LeagueService interface {
	CreateLeague(league *models.League) error
	GetLeagueByID(id uint) (*models.League, error)
	UpdateLeague(league *models.League) error
	DeleteLeague(id uint) error
	GetAllLeagues() ([]*models.League, error)
	GetLeaguesByTeamID(teamID uint) ([]*models.League, error)
	AddTeamToLeague(leagueID, teamID uint) error
	RemoveTeamFromLeague(leagueID, teamID uint) error
	AdvanceWeek(leagueID uint) error
	ViewMatchResults(leagueID uint) ([]*models.Match, error)
	EditMatchResults(matchID uint, updatedMatch *models.Match) error
	PredictChampion(leagueID uint) ([]*dto.TeamPrediction, error)
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

	if league.CurrentWeek > 37 { // TODO write a function inside league entity instead
		return errors.New("league has already ended")
	}

	if len(league.Teams) != 4 {
		return errors.New("league must have exactly 4 teams to advance")
	}

	// Advance the league week
	league, err = s.advanceLeague(league)
	if err != nil {
		return err
	}
	league.CurrentWeek++

	return s.leagueRepo.UpdateLeague(league)
}

// ViewMatchResults returns the match results for the current week
func (s *LeagueServiceImpl) ViewMatchResults(leagueID uint) ([]*models.Match, error) {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return nil, err
	}

	if !league.IsActive() {
		return nil, errors.New("league is not active or has ended")
	}

	matches, err := s.matchRepo.GetMatchesByWeek(leagueID, league.CurrentWeek)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// EditMatchResults allows the user to edit the match results and update the standings accordingly
func (s *LeagueServiceImpl) EditMatchResults(matchID uint, updatedMatch *models.Match) error {
	// Get the existing match
	existingMatch, err := s.matchRepo.GetMatchByID(matchID)
	if err != nil {
		return err
	}

	// Update the match result
	existingMatch.HomeTeamScore = updatedMatch.HomeTeamScore
	existingMatch.AwayTeamScore = updatedMatch.AwayTeamScore

	if err := s.matchRepo.UpdateMatch(existingMatch); err != nil {
		return err
	}

	// Adjust standings
	if err := s.updateTeamStandings(updatedMatch.LeagueID, existingMatch, updatedMatch); err != nil {
		return err
	}

	return nil
}

func (s *LeagueServiceImpl) PredictChampion(leagueID uint) ([]*dto.TeamPrediction, error) {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return nil, err
	}

	if !league.IsActive() {
		return nil, errors.New("league is not active or has ended")
	}

	if league.CurrentWeek < 4 {
		return nil, errors.New("league did not reach the 4th week yet")
	}

	standings := league.Standings
	if len(standings) == 0 {
		return nil, errors.New("no standings found for the league")
	}

	teams := league.Teams
	if len(teams) != 4 {
		return nil, errors.New("league must have 4 teams")
	}

	teamStandings, err := s.combineTeamsAndStandings(teams, standings)
	if err != nil {
		return nil, err
	}

	predictions, err := s.calculateWinProbabilities(teamStandings)
	if err != nil {
		return nil, err
	}

	return predictions, nil
}

func (s *LeagueServiceImpl) PlayAllMatches(leagueID uint) error {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return err
	}

	if !league.IsActive() {
		return errors.New("league is not active or has ended")
	}

	if len(league.Teams) != 4 {
		return errors.New("league must have exactly 4 teams to play matches")
	}

	for league.CurrentWeek < 38 { // TODO refactor
		league, err = s.advanceLeague(league)
		if err != nil {
			return err
		}
		league.CurrentWeek++
	}

	return s.leagueRepo.UpdateLeague(league)
}

// Below are helper functions for simulating matches and calculating scores

func (s *LeagueServiceImpl) advanceLeague(league *models.League) (*models.League, error) {
	// Play matches for the current week
	matches, err := s.playMatches(league)
	if err != nil {
		return nil, err
	}

	// Save match results
	for _, match := range matches {
		err := s.saveMatchResult(&match)
		if err != nil {
			return nil, err
		}
	}

	return league, nil
}

func (s *LeagueServiceImpl) combineTeamsAndStandings(teams []models.Team, standings []models.Standing) ([]teamStanding, error) {
	if len(standings) != 4 {
		return nil, errors.New("4 standings must be present for the league")
	}

	var teamStandings []teamStanding

	// Combine teams with their standings O(1) since each array is guaranteed to have 4 elements
	for _, standing := range standings {
		// search teams array match their id
		for _, team := range teams {
			if team.ID == standing.TeamID {
				teamStandings = append(teamStandings, teamStanding{Team: team, Standing: standing})
			}
		}
	}

	return teamStandings, nil
}

func (s *LeagueServiceImpl) calculateWinProbabilities(teamStandings []teamStanding) ([]*dto.TeamPrediction, error) {
	// Calculate total points, attack strength, and defense strength for all teams
	totalPoints, totalAttackStrength, totalDefenseStrength := s.calculateTotals(teamStandings)

	// Calculate win probabilities
	var predictions []*dto.TeamPrediction
	for _, currentTeamStanding := range teamStandings {
		team := currentTeamStanding.Team
		standing := currentTeamStanding.Standing

		pointsFactor := float64(standing.Points) / float64(totalPoints)
		attackFactor := float64(team.AttackStrength) / float64(totalAttackStrength)
		defenseFactor := float64(team.DefenseStrength) / float64(totalDefenseStrength)

		// Combine factors
		score := pointsFactor*0.5 + attackFactor*0.3 + defenseFactor*0.2

		predictions = append(predictions, &dto.TeamPrediction{
			TeamID:         team.ID,
			TeamName:       team.Name,
			WinProbability: score,
		})
	}

	// Normalize probabilities to sum up to 1
	s.normalizeProbabilities(predictions)

	// Sort predictions by win probability
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].WinProbability > predictions[j].WinProbability
	})

	return predictions, nil
}

func (s *LeagueServiceImpl) calculateTotals(teamStandings []teamStanding) (totalPoints, totalAttackStrength, totalDefenseStrength int) {
	for _, currentTeamStanding := range teamStandings {
		standing := currentTeamStanding.Standing
		team := currentTeamStanding.Team

		totalPoints += standing.Points
		totalAttackStrength += team.AttackStrength
		totalDefenseStrength += team.DefenseStrength
	}
	return
}

func (s *LeagueServiceImpl) normalizeProbabilities(predictions []*dto.TeamPrediction) {
	totalScore := 0.0
	for _, prediction := range predictions {
		totalScore += prediction.WinProbability
	}

	for i := range predictions {
		predictions[i].WinProbability /= totalScore
	}
}

// Custom object to store teams with their standings
type teamStanding struct {
	models.Team
	models.Standing
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

// saveMatchResult saves the match result and updates the standings
func (s *LeagueServiceImpl) saveMatchResult(match *models.Match) error {
	if err := s.matchRepo.CreateMatch(match); err != nil {
		return err
	}

	return s.updateTeamStandings(match.LeagueID, nil, match)
}

// updateTeamStandings updates the standings based on old and new match results for both home and away teams
func (s *LeagueServiceImpl) updateTeamStandings(leagueID uint, oldMatch, newMatch *models.Match) error {
	// Revert old match results if oldMatch is not nil
	if oldMatch != nil {
		if err := s.adjustStandings(leagueID, oldMatch.HomeTeamID, oldMatch.HomeTeamScore, oldMatch.AwayTeamScore, true); err != nil {
			return err
		}
		if err := s.adjustStandings(leagueID, oldMatch.AwayTeamID, oldMatch.AwayTeamScore, oldMatch.HomeTeamScore, true); err != nil {
			return err
		}
	}

	// Apply new match results
	if newMatch != nil {
		if err := s.adjustStandings(leagueID, newMatch.HomeTeamID, newMatch.HomeTeamScore, newMatch.AwayTeamScore, false); err != nil {
			return err
		}
		if err := s.adjustStandings(leagueID, newMatch.AwayTeamID, newMatch.AwayTeamScore, newMatch.HomeTeamScore, false); err != nil {
			return err
		}
	}

	return nil
}

// adjustStandings adjusts the standings for a team based on match results
func (s *LeagueServiceImpl) adjustStandings(leagueID, teamID uint, teamScore, opponentScore int, isRevert bool) error {
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

	// Revert old match result if needed
	if isRevert {
		standings.GoalDifference -= teamScore - opponentScore
		standings.Played--

		if teamScore > opponentScore {
			standings.Wins--
			standings.Points -= 3
		} else if teamScore == opponentScore {
			standings.Draws--
			standings.Points--
		} else {
			standings.Losses--
		}
	} else {
		// Apply new match result
		standings.GoalDifference += teamScore - opponentScore
		standings.Played++

		if teamScore > opponentScore {
			standings.Wins++
			standings.Points += 3
		} else if teamScore == opponentScore {
			standings.Draws++
			standings.Points++
		} else {
			standings.Losses++
		}
	}

	// Standing is newly created if err is not nil
	if err != nil {
		return s.standingRepo.CreateStanding(standings)
	}
	return s.standingRepo.UpdateStanding(standings)
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
