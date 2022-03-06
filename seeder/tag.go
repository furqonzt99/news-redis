package seeder

import (
	"fmt"

	"github.com/furqonzt99/news-redis/domain/entity"
	"gorm.io/gorm"
)

func TagSeeder(db *gorm.DB) {
	for i := 1; i <= 10; i++ {
		db.Create(&entity.Tag{
			Name: "Topic" + fmt.Sprint(i),
		})
	}
}
