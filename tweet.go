package main

import "gorm.io/gorm"

type tweet struct {
	gorm.Model
	UserId  int64  `json:"userid"`
	Content string `json:"content"`
}
