package handler

import (
	"API-CRYPT/src/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type KlinesHandler struct {
	service *services.KlinesService
}

func NewKlinesHandler(service *services.KlinesService) *KlinesHandler {
	return &KlinesHandler{service: service}
}

func (h *KlinesHandler) GetKlines(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	interval := r.URL.Query().Get("interval")
	daysStr := r.URL.Query().Get("days")

	if symbol == "" || interval == "" || daysStr == "" {
		http.Error(w, "Необходимо указать symbol, interval и days", http.StatusBadRequest)
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		http.Error(w, "days должен быть числом", http.StatusBadRequest)
		return
	}

	candles, err := h.service.GetKlines(symbol, interval, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formattedJSON, err := json.MarshalIndent(candles, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка при преобразовании в JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(formattedJSON)
}

func (h *KlinesHandler) GetAvailableSymbols(w http.ResponseWriter, r *http.Request) {
	symbols, err := h.service.GetAvailableSymbols()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.MarshalIndent(symbols, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка при преобразовании в JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *KlinesHandler) GetKlinesCSV(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	interval := r.URL.Query().Get("interval")
	daysStr := r.URL.Query().Get("days")

	if symbol == "" || interval == "" || daysStr == "" {
		http.Error(w, "Необходимо указать symbol, interval и days", http.StatusBadRequest)
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		http.Error(w, "days должен быть числом", http.StatusBadRequest)
		return
	}

	fileName, err := h.service.CreateCSVFile(symbol, interval, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("http://localhost:8080/static/%s", fileName)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"file_url": "%s"}`, fileURL)))
}
