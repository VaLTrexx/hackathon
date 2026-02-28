package models

type Event struct {
	Actor1Name     string  `json:"actor1_name"`
	Actor2Name     string  `json:"actor2_name"`
	EventCode      string  `json:"event_code"`
	GoldsteinScale float64 `json:"goldstein_scale"`
	NumMentions    int     `json:"num_mentions"`
	AvgTone        float64 `json:"avg_tone"`
	EventDate      string  `json:"event_date"`
}
