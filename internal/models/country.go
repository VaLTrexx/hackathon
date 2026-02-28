package models

type Country struct {
	Name       string `json:"name"`
	CCA2       string `json:"cca2"`
	CCA3       string `json:"cca3"`
	Population int64  `json:"population"`
	Region     string `json:"region"`
	Currency   string `json:"currency"`
}
