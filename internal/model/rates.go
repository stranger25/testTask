package model

import "encoding/xml"

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