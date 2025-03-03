package services

import (
	"API-CRYPT/src/payload"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	bybitAPIURL = "https://api.bybit.com/v5/market/kline"
)

type KlinesService struct{}

func NewKlinesService() *KlinesService {
	return &KlinesService{}
}

func (s *KlinesService) GetKlinesWithIntervals(symbol, interval string, days int) ([]payload.FormattedCandle, error) {
	var allCandles []payload.FormattedCandle

	endTime := time.Now().Unix() * 1000
	startTime := time.Now().Add(-time.Duration(days)*24*time.Hour).Unix() * 1000

	intervalDuration := 12 * time.Hour
	currentEndTime := endTime
	for currentEndTime > startTime {
		currentStartTime := currentEndTime - int64(intervalDuration.Seconds())*1000
		if currentStartTime < startTime {
			currentStartTime = startTime
		}

		url := fmt.Sprintf("%s?category=spot&symbol=%s&interval=%s&start=%d&end=%d", bybitAPIURL, symbol, interval, currentStartTime, currentEndTime)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ошибка при выполнении запроса к Bybit API: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении ответа от Bybit API: %v", err)
		}

		var result payload.BybitResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("ошибка при парсинге JSON: %v", err)
		}

		for _, candle := range result.Result.List {
			var openTime int64
			switch v := candle[0].(type) {
			case float64:
				openTime = int64(v)
			case string:
				openTime, _ = strconv.ParseInt(v, 10, 64)
			default:
				return nil, fmt.Errorf("неподдерживаемый тип данных для open_time")
			}

			openTimeFormatted := time.Unix(openTime/1000, 0).Format("2006-01-02 15:04:05")

			allCandles = append(allCandles, payload.FormattedCandle{
				OpenTime:    openTimeFormatted,
				OpenPrice:   fmt.Sprint(candle[1]),
				HighPrice:   fmt.Sprint(candle[2]),
				LowPrice:    fmt.Sprint(candle[3]),
				ClosePrice:  fmt.Sprint(candle[4]),
				VolumeBase:  fmt.Sprint(candle[5]),
				VolumeQuote: fmt.Sprint(candle[6]),
			})
		}

		currentEndTime = currentStartTime - 1
	}

	return allCandles, nil
}

func (s *KlinesService) GetAvailableSymbols() ([]string, error) {
	url := "https://api.bybit.com/v5/market/instruments-info?category=spot"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса к Bybit API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа от Bybit API: %v", err)
	}

	var result struct {
		Result struct {
			List []struct {
				Symbol string `json:"symbol"`
			} `json:"list"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге JSON: %v", err)
	}

	var symbols []string
	for _, instrument := range result.Result.List {
		symbols = append(symbols, instrument.Symbol)
	}

	return symbols, nil
}

func (s *KlinesService) CreateCSVFile(symbol, interval string, days int) (string, error) {
	candles, err := s.GetKlinesWithIntervals(symbol, interval, days)
	if err != nil {
		return "", err
	}

	if len(candles) == 0 {
		return "", &payload.NoDataError{Symbol: symbol}
	}

	fileName := fmt.Sprintf("%s_%s_%dd.csv", symbol, interval, days)
	filePath := "../static/" + fileName

	if _, err := os.Stat("../static"); os.IsNotExist(err) {
		if err := os.Mkdir("../static", 0755); err != nil {
			return "", fmt.Errorf("ошибка при создании папки static: %v", err)
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Open Time", "Open Price", "High Price", "Low Price", "Close Price", "Volume Base", "Volume Quote"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("ошибка при записи заголовков CSV: %v", err)
	}

	for _, candle := range candles {
		record := []string{
			candle.OpenTime,
			candle.OpenPrice,
			candle.HighPrice,
			candle.LowPrice,
			candle.ClosePrice,
			candle.VolumeBase,
			candle.VolumeQuote,
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("ошибка при записи данных в CSV: %v", err)
		}
	}

	return fileName, nil
}
