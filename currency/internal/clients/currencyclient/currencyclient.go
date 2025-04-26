package currencyclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lekss361/currencyservice/currency/internal/dto"
)

func GetCurs() (*dto.RatesResponse, error) {
	resp, err := http.Get("https://latest.currency-api.pages.dev/v1/currencies/rub.json")
	if err != nil {
		return nil, fmt.Errorf("http.Get error: %w", err)
	}
	defer resp.Body.Close()

	var tmp struct {
		Date string             `json:"date"`
		Rub  map[string]float64 `json:"rub"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return nil, fmt.Errorf("json.Decode error: %w", err)
	}

	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, tmp.Date)
	if err != nil {
		return nil, fmt.Errorf("time.Parse error: %w", err)
	}

	result := &dto.RatesResponse{
		Date: parsedDate,
		Rub:  tmp.Rub,
	}
	return result, nil
}
