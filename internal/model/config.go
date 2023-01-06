package model

type Config struct {
	Database struct {
		Connection string `json:"connection"`
	} `json:"database"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}
