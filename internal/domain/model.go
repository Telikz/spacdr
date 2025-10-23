package domain

import "time"

type Card struct {
	Front      string    `json:"front"`
	Back       string    `json:"back"`
	Score      int       `json:"score"`
	LastReview time.Time `json:"last_review"`
}

type Deck struct {
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}
