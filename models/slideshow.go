package models

import "gorm.io/gorm"

type Slideshow struct {
	gorm.Model
	Url string `json:"url"`
	Title string `json:"title"`
}
