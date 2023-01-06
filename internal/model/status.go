package model

type Success struct {
	Success bool `json:"success"`
}

type UnSuccess struct {
	Errmsg  string `json:"errmsg"`
}