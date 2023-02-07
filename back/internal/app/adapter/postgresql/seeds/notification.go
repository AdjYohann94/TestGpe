package seeds

import (
	"fmt"
	"math/rand"
	"time"

	"gpe_project/internal/app/adapter/postgresql/model"
)

func GetNotificationSeedForUser(userID uint, number int) []*model.Notification {
	notifications := make([]*model.Notification, number)
	for i := 0; i < number; i++ {
		notifications[i] = &model.Notification{
			UserID:  userID,
			Read:    false,
			Message: fmt.Sprintf("notification %d", i),
			Level:   1,
		}
	}

	return notifications
}

func GetBarometers(userId uint, number int) []*model.Barometer {
	barometers := make([]*model.Barometer, number)
	types := [4]string{"happy", "work", "bar", "foo"}
	for i := 0; i < number; i++ {
		d := time.Now().AddDate(0, 0, i)
		barometers[i] = &model.Barometer{
			UserID: userId,
			Score:  rand.Intn(100),
			Date:   &d,
			Type:   types[rand.Intn(len(types))],
		}
	}

	return barometers
}
