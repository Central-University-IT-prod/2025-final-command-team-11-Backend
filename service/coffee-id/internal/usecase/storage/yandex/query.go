package yandex

import (
	"fmt"
	"time"

	"github.com/nikitaSstepanov/coffee-id/internal/entity"
)

func idQuery(id string) string {
	return fmt.Sprintf(
		`
			SELECT * FROM %s 
			WHERE yandex_id = '%s';
		`, yandexTable, id,
	)
}

func userIdQuery(userId uint64) string {
	return fmt.Sprintf(
		`
			SELECT * FROM %s 
			WHERE user_id = %d;
		`, yandexTable, userId,
	)
}

func createQuery(yndx *entity.Yandex) string {
	return fmt.Sprintf(
		`
			INSERT INTO %s 
				(yandex_id, email, name, birthday, user_id) 
			VALUES 
				('%s', '%s', '%s', '%s', %d);
		`, yandexTable, yndx.YandexId, yndx.Email, yndx.Name, yndx.Birthday.Format(time.DateOnly), yndx.UserId,
	)
}

func redisKey(userId uint64) string {
	return fmt.Sprintf("yandex:%d", userId)
}
