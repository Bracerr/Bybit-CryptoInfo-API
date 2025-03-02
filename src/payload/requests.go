package payload

type KlinesRequest struct {
	Symbol   string `json:"symbol"`
	Interval string `json:"interval"`
	Days     int    `json:"days"`
}
