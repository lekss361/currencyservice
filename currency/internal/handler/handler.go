package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lekss361/currencyservice/currency/internal/dto"
)

// Handler handles HTTP requests for currency data.
type Handler struct {
	repo dto.RatesRepo
}

// New creates a new Handler with the given RatesRepo.
func New(repo dto.RatesRepo) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes registers HTTP routes for currency data.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/rate", h.GetRateByDate)
	mux.HandleFunc("/history", h.GetHistory)
}

// GetRateByDate handles GET /rate?date=YYYY-MM-DD
func (h *Handler) GetRateByDate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "missing date parameter", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid date format: %v", err), http.StatusBadRequest)
		return
	}
	rates, err := h.repo.Get(date)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching rates: %v", err), http.StatusInternalServerError)
		return
	}
	resp := dto.RatesResponse{
		Date: date,
		Rub:  rates,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetHistory handles GET /history?start=YYYY-MM-DD&end=YYYY-MM-DD
func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	if startStr == "" || endStr == "" {
		http.Error(w, "missing start or end parameter", http.StatusBadRequest)
		return
	}
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid start date format: %v", err), http.StatusBadRequest)
		return
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid end date format: %v", err), http.StatusBadRequest)
		return
	}
	if end.Before(start) {
		http.Error(w, "end date must not be before start date", http.StatusBadRequest)
		return
	}

	type DateRates struct {
		Date string             `json:"date"`
		Rub  map[string]float64 `json:"rub"`
	}
	var history []DateRates
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		rates, err := h.repo.Get(d)
		if err != nil {
			http.Error(w, fmt.Sprintf("error fetching rates for %s: %v", d.Format("2006-01-02"), err), http.StatusInternalServerError)
			return
		}
		history = append(history, DateRates{
			Date: d.Format("2006-01-02"),
			Rub:  rates,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
