package repository

import (
	"github.com/furqonzt99/news-redis/domain/entity"
	"gorm.io/gorm"
)

type TagInterface interface {
	Create(tag entity.Tag) (entity.Tag, error)
	ReadAll() ([]entity.Tag, error)
	Edit(id int, newTag entity.Tag) (entity.Tag, error)
	Delete(id int) (entity.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *tagRepository {
	return &tagRepository{db: db}
}

func (tr *tagRepository) Create(tag entity.Tag) (entity.Tag, error) {
	if err := tr.db.Create(&tag).Error; err != nil {
		return tag, err
	}

	return tag, nil
}

func (tr *tagRepository) ReadAll() ([]entity.Tag, error) {
	var tags []entity.Tag

	tr.db.Find(&tags)

	return tags, nil
}

func (tr *tagRepository) Edit(id int, newTag entity.Tag) (entity.Tag, error) {
	var tag entity.Tag

	if err := tr.db.First(&tag, id).Error; err != nil {
		return tag, err
	}

	tr.db.Model(&tag).Updates(newTag)

	return tag, nil
}

func (tr *tagRepository) Delete(id int) (entity.Tag, error) {
	var tag entity.Tag

	if err := tr.db.First(&tag, id).Error; err != nil {
		return tag, err
	}

	tr.db.Delete(&tag)

	return tag, nil
}
