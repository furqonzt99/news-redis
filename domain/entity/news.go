package entity

import (
	"errors"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title  string
	Body   string
	Status string `gorm:"default:draft"`
	Tags   []Tag  `gorm:"many2many:news_tags;"`
}

type NewsFilter struct {
	Status string
	Tags   []string
}

type NewsTags struct {
	NewsID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

func (NewsTags) BeforeCreate(db *gorm.DB) error {
	err := db.SetupJoinTable(&News{}, "Tags", &NewsTags{})

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
