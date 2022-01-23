package models

import "gorm.io/gorm"

type Announcement struct {
	gorm.Model
	Title 		string
	Content 	string
	PictureUrl 	string
	AdminID 	uint
}
