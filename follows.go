package main

import "gorm.io/gorm"

type follows struct {
	gorm.Model
	SourceUser string `json:"sourceuser"`
	TargetUser string `json:"targetuser"`
}
