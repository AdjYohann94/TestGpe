package repository

import (
	"gorm.io/gorm"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/postgresql/scope"
)

type NotificationRepository struct {
	DB *gorm.DB
}

// CreateNotification add a user notification with the default model.LevelMinimum priority.
func (nr NotificationRepository) CreateNotification(notification model.Notification) (*model.Notification, error) {
	result := nr.DB.Create(&notification)
	err := CatchCreateError(result.Error, result.RowsAffected)
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

// RetrieveUnreadNotifications returns all unread notifications for the userID specified. Notifications are sorted by
// creationDate reversed.
func (nr NotificationRepository) RetrieveUnreadNotifications(userID uint, filters scope.Filters) ([]*model.Notification, error) {
	var notifications []*model.Notification
	result := nr.DB.
		Scopes(scope.AutoFilterScore(filters)).
		Where(&model.Notification{UserID: userID}).
		Where("read = ?", false).
		Order("created_at desc").
		Find(&notifications)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// RetrieveAllNotifications returns all notifications for the userID specified. Notifications are sorted by
// creationDate reversed.
func (nr NotificationRepository) RetrieveAllNotifications(userID uint, filters scope.Filters) ([]*model.Notification, error) {
	var notifications []*model.Notification
	result := nr.DB.
		Scopes(scope.AutoFilterScore(filters)).
		Where(&model.Notification{UserID: userID}).
		Order("created_at desc").
		Find(&notifications)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// ReadNotifications set all notifications for the user to status read true.
func (nr NotificationRepository) ReadNotifications(userID uint) error {
	result := nr.DB.
		Model(&model.Notification{}).
		Where("user_id = ? and read = ?", userID, false).
		Update("read", "true")
	err := CatchUpdateError(result.Error, result.RowsAffected)
	if err != nil {
		return err
	}

	return nil
}
