package models

type EventInput struct {
	Tone           float64 `json:"tone"`
	Mentions       int     `json:"mentions"`
	SeverityWeight float64 `json:"severity_weight"`
	Date           string  `json:"date,omitempty"`
}

type RiskRequest struct {
	Events []EventInput `json:"events"`
	Prices []float64    `json:"prices"`
}