package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ValTrexx/cs/internal/models"
)

type alphaResponse struct {
	GlobalQuote struct {
		Symbol string `json:"01. symbol"`
		Open   string `json:"02. open"`
		High   string `json:"03. high"`
		Low    string `json:"04. low"`
		Price  string `json:"05. price"`
		Volume string `json:"06. volume"`
	} `json:"Global Quote"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchMarket(symbol string) (models.Market, error) {

	apiKey := os.Getenv("PVJTUC3J4HTITE9P")

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol,
		apiKey,
	)

	resp, err := client.Get(url)
	if err != nil {
		return models.Market{}, err
	}
	defer resp.Body.Close()

	var result alphaResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return models.Market{}, err
	}

	open, _ := strconv.ParseFloat(result.GlobalQuote.Open, 64)
	high, _ := strconv.ParseFloat(result.GlobalQuote.High, 64)
	low, _ := strconv.ParseFloat(result.GlobalQuote.Low, 64)
	closePrice, _ := strconv.ParseFloat(result.GlobalQuote.Price, 64)
	volume, _ := strconv.ParseInt(result.GlobalQuote.Volume, 10, 64)

	return models.Market{
		Symbol: result.GlobalQuote.Symbol,
		Open:   open,
		High:   high,
		Low:    low,
		Close:  closePrice,
		Volume: volume,
	}, nil
}
