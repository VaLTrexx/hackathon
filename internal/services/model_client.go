package services

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ValTrexx/hackathon/internal/models"
)

func SendToRiskModel(req models.RiskRequest) (map[string]interface{}, error) {

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		"http://localhost:8000/compute_risk",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}