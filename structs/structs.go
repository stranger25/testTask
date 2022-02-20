package structs

import "encoding/xml"

type Config struct {
	Database struct {
		Connection string `json:"connection"`
	} `json:"database"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

type Rates struct {
	XMLName     xml.Name `xml:"rates"`
	Text        string   `xml:",chardata"`
	Generator   string   `xml:"generator"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Copyright   string   `xml:"copyright"`
	Date        string   `xml:"date"`
	Item        []struct {
		Text        string `xml:",chardata"`
		Fullname    string `xml:"fullname"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Quant       string `xml:"quant"`
		Index       string `xml:"index"`
		Change      string `xml:"change"`
	} `xml:"item"`
}

type Success struct {
	Success bool `json:"success"`
}

type UnSuccess struct {
	Errmsg  string `json:"errmsg"`
}

type RatesOut struct {
	Title string  `json:"title"`
	Code  string  `json:"code"`
	Value float64 `json:"value"`
	Date  string  `json:"date"`
}
