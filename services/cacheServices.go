package services

import (
	"context"
	"fmt"

	"github.com/furqonzt99/news-redis/constants"
)

var ctx = context.Background()

func CreateCache(entity string, id int, filter interface{}, data interface{}) error {
	key := entity + ":" + fmt.Sprint(id) + ":" + fmt.Sprint(filter)

	err := constants.Rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCache(entity string, id int, filter interface{}) (string, error) {
	key := entity + ":" + fmt.Sprint(id) + ":" + fmt.Sprint(filter)

	data, err := constants.Rdb.Get(ctx, key).Result()
	if err != nil {
		return data, err
	}

	return data, nil
}

func DeleteCache(entity string) error {
	iter := constants.Rdb.Scan(ctx, 0, entity+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := constants.Rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
