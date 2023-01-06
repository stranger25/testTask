package model

type RatesOut struct {
	Title string  `json:"title"`
	Code  string  `json:"code"`
	Value float64 `json:"value"`
	Date  string  `json:"date"`
}
