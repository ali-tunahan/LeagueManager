package models

import "gorm.io/gorm"

// Team represents a Football team, may be affiliated with multiple leagues.
type Team struct {
	gorm.Model
	Name            string `json:"name"`
	AttackStrength  int    `json:"attack_strength"`
	DefenseStrength int    `json:"defense_strength"`
}
