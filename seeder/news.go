package seeder

import (
	"math/rand"

	"github.com/furqonzt99/news-redis/domain/entity"
	"gorm.io/gorm"
)

func NewsSeeder(db *gorm.DB) {
	status := []string{"draft", "publish", "deleted"}
	for i := 1; i <= 100; i++ {
		db.Create(&entity.News{
			Title:  "Title",
			Body:   "Body",
			Status: status[rand.Intn(3-0)],
		})
	}
}

func NewsTagsSeeder(db *gorm.DB) {
	for i := 1; i <= 300; i++ {
		db.Create(&entity.NewsTags{
			NewsID: uint(rand.Intn(100-1) + 1),
			TagID:  uint(rand.Intn(10-1) + 1),
		})
	}
}
