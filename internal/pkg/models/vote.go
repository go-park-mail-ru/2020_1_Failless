package models

import "time"

type Vote struct {
	Uid   int       `json:"uid"`
	Id    int       `json:"id"`
	Value int8      `json:"value"`
	Date  time.Time `json:"-"`
}

