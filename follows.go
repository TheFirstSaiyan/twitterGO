package main

import "gorm.io/gorm"

type follows struct {
	gorm.Model
	SourceId int64 `json:"sourceid"`
	TargetId int64 `json:"targetid"`
}
