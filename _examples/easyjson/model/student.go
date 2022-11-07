package model

import "time"

// //go:generate easyjson -all student.go
//
//easyjson:json
type School struct {
	Name string `json:"school_name"`
	Addr string `json:"school_addr"`
}

//easyjson:json
type Student struct {
	Id       int       `json:"id"`
	Name     string    `json:"student_name"`
	School   School    `json:"student_school"`
	Birthday time.Time `json:"birthday"`
}
