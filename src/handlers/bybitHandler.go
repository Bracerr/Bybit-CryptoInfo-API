package handler

import (
	"API-CRYPT/src/payload"
	"API-CRYPT/src/services"
	"encoding/json"
	"errors"
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
		h.respondWithError(w, http.StatusBadRequest, "Необходимо указать symbol, interval и days")
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "days должен быть числом")
		return
	}

	candles, err := h.service.GetKlinesWithIntervals(symbol, interval, days)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if candles == nil {
		h.respondWithError(w, http.StatusBadRequest, "Нет данных для данной валюты: "+symbol)
		return
	}

	formattedJSON, err := json.MarshalIndent(candles, "", "  ")
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Ошибка при преобразовании в JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(formattedJSON)
}

func (h *KlinesHandler) GetAvailableSymbols(w http.ResponseWriter, r *http.Request) {
	symbols, err := h.service.GetAvailableSymbols()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseJSON, err := json.MarshalIndent(symbols, "", "  ")
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Ошибка при преобразовании в JSON")
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
		h.respondWithError(w, http.StatusBadRequest, "Необходимо указать symbol, interval и days")
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "days должен быть числом")
		return
	}

	fileName, err := h.service.CreateCSVFile(symbol, interval, days)
	if err != nil {
		var noDataError *payload.NoDataError
		if errors.As(err, &noDataError) {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fileURL := fmt.Sprintf("http://localhost:8080/static/%s", fileName)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"file_url": "%s"}`, fileURL)))
}

func (h *KlinesHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := payload.ErrorResponse{Message: message}
	jsonResponse, _ := json.Marshal(errorResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
