package model

type Message struct {
	Type  string `json:"type"`
	Data  string `json:"data"`
	Model string `json:"model"`
}
