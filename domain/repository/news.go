package repository

import (
	"github.com/furqonzt99/news-redis/domain/entity"
	"gorm.io/gorm"
)

type NewsInterface interface {
	Create(news entity.News, tags []int) (entity.News, error)
	ReadAll(filter entity.NewsFilter) ([]entity.News, error)
	ReadOne(id int) (entity.News, error)
	Edit(id int, newNews entity.News, tags []int) (entity.News, error)
	Delete(id int) (entity.News, error)
	SetStatusDeleted(id int) (entity.News, error)
	SetStatusPublish(id int) (entity.News, error)
	SetStatusDraft(id int) (entity.News, error)
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *newsRepository {
	return &newsRepository{db: db}
}

func (nr *newsRepository) Create(news entity.News, tags []int) (entity.News, error) {
	if err := nr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&news).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			if err := tx.Create(&entity.NewsTags{
				NewsID: news.ID,
				TagID:  uint(tag),
			}).Error; err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return news, err
	}

	return news, nil
}

func (nr *newsRepository) ReadAll(filter entity.NewsFilter) ([]entity.News, error) {
	var news []entity.News

	if filter.Tags[0] != "" && filter.Status != "" {
		nr.db.Preload("Tags", "name IN ?", filter.Tags).Where("status = ?", filter.Status).Find(&news)
	} else if filter.Tags[0] != "" {
		nr.db.Preload("Tags", "name IN ?", filter.Tags).Find(&news)
	} else if filter.Status != "" {
		nr.db.Preload("Tags").Where("status = ?", filter.Status).Find(&news)
	} else {
		nr.db.Preload("Tags").Find(&news)
	}

	return news, nil
}

func (nr *newsRepository) ReadOne(id int) (entity.News, error) {
	var news entity.News

	if err := nr.db.Preload("Tags").First(&news, id).Error; err != nil {
		return news, err
	}

	return news, nil
}

func (nr *newsRepository) Edit(id int, newNews entity.News, tags []int) (entity.News, error) {
	var news entity.News

	if err := nr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&news, id).Error; err != nil {
			return err
		}

		nr.db.Model(&news).Updates(newNews)

		if err := tx.Delete(&entity.NewsTags{}, "news_id = ?", news.ID).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			if err := tx.Create(&entity.NewsTags{
				NewsID: news.ID,
				TagID:  uint(tag),
			}).Error; err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return news, err
	}

	return news, nil
}

func (nr *newsRepository) Delete(id int) (entity.News, error) {
	var news entity.News

	if err := nr.db.First(&news, id).Error; err != nil {
		return news, err
	}

	nr.db.Delete(&news, id)

	return news, nil
}

func (nr *newsRepository) SetStatusDeleted(id int) (entity.News, error) {
	var news entity.News

	if err := nr.db.First(&news, id).Error; err != nil {
		return news, err
	}

	nr.db.Model(&news).Update("status", "deleted")

	return news, nil
}

func (nr *newsRepository) SetStatusPublish(id int) (entity.News, error) {
	var news entity.News

	if err := nr.db.First(&news, id).Error; err != nil {
		return news, err
	}

	nr.db.Model(&news).Update("status", "publish")

	return news, nil
}

func (nr *newsRepository) SetStatusDraft(id int) (entity.News, error) {
	var news entity.News

	if err := nr.db.First(&news, id).Error; err != nil {
		return news, err
	}

	nr.db.Model(&news).Update("status", "draft")

	return news, nil
}
