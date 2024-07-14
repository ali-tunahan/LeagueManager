# LeagueManager Documentation

## Table of Contents
1. [Introduction](#introduction)
2. [Business Rules](#business-rules)
3. [API Endpoints](#api-endpoints)
4. [Getting Started](#getting-started)
5. [How to Use](#how-to-use)
6. [Running Tests](#running-tests)

## Introduction
LeagueManager is a robust league management system designed to handle the complexities of managing sports leagues. This system allows for the creation of leagues, management of teams, scheduling of matches, and prediction of league champions. This documentation serves to provide detailed information about the business rules, API usage, and setup instructions for developers.

## Business Rules
1. **League Creation**: A league can be created with a name. Leagues are created with no teams initially. Teams can be added to the league later. A league starts with week 0, indicating that it has not started yet.
2. **Starting a League**: A league must be started before any matches can be played. Once started, the league week advances from 0 to 1.
3. **Team Management**: Teams can be added to or removed from leagues. Each team has attributes like name, attack strength, and defense strength. A team can belong to multiple leagues.
4. **Team Removal**: If the league has started, teams cannot be removed from the league. Teams can only be removed before the league starts.
5. **Match Scheduling**: Matches are scheduled automatically when a league is started. Each team plays every other team twice (home and away).
6. **League Advancement**: Leagues advance week by week. Each week, scheduled matches are played, and results are recorded. When the league is at week 1, the matches for week 1 will be played when advanced. After advancing, the week is incremented (e.g., from 1 to 2). So, the week count indicates the week of the league that was not played yet.
7. **Match Results**: Match results can be viewed, and match details can be edited if necessary.
8. **Champion Prediction**: The system can predict the champion based on current standings and match results.
9. **End of Season**: A league season consists of 38 weeks. At the end of the season, the league champion is determined based on standings. Week 0 means has not started and week 39 means league is completed.
10. **Initialization for Testing**: A special function can initialize a league with predefined teams (e.g., Premier League teams).

## API Endpoints

### Team Endpoints
- **POST /api/teams**: Add a new team.
- **GET /api/teams**: Get all teams.
- **GET /api/teams/:teamID**: Get a team by ID.
- **PUT /api/teams/:teamID**: Update a team.
- **DELETE /api/teams/:teamID**: Delete a team.

### League Endpoints
- **POST /api/leagues/create**: Create a new league.
- **POST /api/leagues/initialize**: Create and initialize a league with default teams.
- **POST /api/leagues/add-team/:leagueID/:teamID**: Add a team to a league.
- **POST /api/leagues/remove-team/:leagueID/:teamID**: Remove a team from a league.
- **POST /api/leagues/start/:leagueID**: Start the league by setting up initial matches.
- **POST /api/leagues/advance-week/:leagueID**: Advance the league by one week.
- **GET /api/leagues/view-matches/:leagueID**: View match results for the current week.
- **POST /api/leagues/edit-match/:matchID**: Edit match results.
- **GET /api/leagues/predict-champion/:leagueID**: Predict the champion of the league.
- **POST /api/leagues/play-all-matches/:leagueID**: Play all remaining matches in the league.

## Getting Started

### Prerequisites
- Go 1.16 or later
- SQLite3
- Gin Framework
- GORM (Go Object Relational Mapping)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/LeagueManager.git
    ```
2. Navigate to the project directory:
    ```sh
    cd LeagueManager
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```
4. Run the application:
    ```sh
    go run main.go
    ```

## How to Use

### Creating a League

To create a new league, send a POST request to `/api/leagues/create` with the league name:
```json
{
   "name": "Test League"
}
```

### Adding Teams

To add a team, send a POST request to `/api/teams` with the team details:
```json
{
  "name": "Team A",
  "attack_strength": 80,
  "defense_strength": 75
}
```

### Starting a League

To start the league and generate the initial match schedule, send a POST request to `/api/leagues/start/:leagueID`.

### Advancing the League

To advance the league by one week and play the matches scheduled for that week, send a POST request to `/api/leagues/advance-week/:leagueID`.

### Viewing Match Results

To view the match results for the current week, send a GET request to `/api/leagues/view-matches/:leagueID`.

### Predicting the Champion

To predict the champion of the league, send a GET request to `/api/leagues/predict-champion/:leagueID`.

## Running Tests

### Prerequisites
- Testify package for assertions
- Gin Framework's testing package

### Running Tests

To run the tests, execute the following command:
```sh
go test ./...
```
This will run all the unit tests and output the results.
Tests are also ran on commits and PRs on the main branch using GitHub Actions.

