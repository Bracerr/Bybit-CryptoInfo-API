package payload

import "fmt"

type FormattedCandle struct {
	OpenTime    string `json:"open_time"`    // Время открытия в читаемом формате
	OpenPrice   string `json:"open_price"`   // Цена открытия
	HighPrice   string `json:"high_price"`   // Максимальная цена
	LowPrice    string `json:"low_price"`    // Минимальная цена
	ClosePrice  string `json:"close_price"`  // Цена закрытия
	VolumeBase  string `json:"volume_base"`  // Объем в базовой валюте (BTC)
	VolumeQuote string `json:"volume_quote"` // Объем в котировочной валюте (USDT)
}

type BybitResponse struct {
	Result struct {
		Category string          `json:"category"`
		List     [][]interface{} `json:"list"` // Список свечей
		Symbol   string          `json:"symbol"`
	} `json:"result"`
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Time    int64  `json:"time"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type NoDataError struct {
	Symbol string
}

func (e *NoDataError) Error() string {
	return fmt.Sprintf("нет данных для валюты: %s", e.Symbol)
}
