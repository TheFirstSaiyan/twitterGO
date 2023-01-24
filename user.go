package main

import "gorm.io/gorm"

type user struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique"`
	Password string `json:"password"`
}
