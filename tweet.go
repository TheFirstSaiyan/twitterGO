package main

import "gorm.io/gorm"

type tweet struct {
	gorm.Model
	UserName string `json:"name"`
	Content  string `json:"content"`
}
