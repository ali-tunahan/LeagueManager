package models

// Team represents a Football team, may be affiliated with multiple leagues.
type Team struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	Name            string `json:"name"`
	AttackStrength  int    `json:"attack_strength"`
	DefenseStrength int    `json:"defense_strength"`
}
