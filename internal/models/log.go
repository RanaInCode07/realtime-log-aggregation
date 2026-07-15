package models

import (
	"time"
)
type LogEvent struct{
	ID          string			`json:"id"`
	TimeStamp   time.Time		`json:"time"`
	Level       string			`json:"level"`
	Service     string			`json:"service"`
	Message     string 			`json:"message"`
}



