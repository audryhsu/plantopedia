package models

import "gorm.io/gorm"

type Plant struct {
	gorm.Model
	Id int `json:"id" gorm:"primaryKey"`
	Species     string `json:"species"`
	Description string `json:"description"`
	Water int `json:"water"`
}