package repository

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func CatchFindError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil
	default:
		log.Println("[ERROR] ", err.Error())
		return fmt.Errorf("something wrong happens")
	}
}

func CatchFindOrFailError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &NotFoundError{}
	default:
		log.Println("[ERROR] ", err.Error())
		return fmt.Errorf("something wrong happens")
	}
}

func CatchCreateError(err error, rowsAffected int64) error {
	switch {
	case err == nil:
		return nil
	case rowsAffected == 0:
		return fmt.Errorf("an error occured when creating the item")
	default:
		log.Println("[ERROR] ", err.Error())
		return fmt.Errorf("something wrong happens")
	}
}

func CatchError(err error) error {
	if err == nil {
		return nil
	}
	log.Println("[ERROR] ", err.Error())
	return fmt.Errorf("something wrong happens")
}

func CatchFindAllError(err error) error {
	if err == nil {
		return nil
	}
	log.Println("[ERROR] ", err.Error())
	return fmt.Errorf("something wrong happens")
}

func CatchDeleteError(err error, rowsAffected int64) error {
	if err == nil {
		if rowsAffected == 0 {
			return &NotFoundError{}
		}
		return nil
	}
	log.Println("[ERROR] ", err.Error())
	return fmt.Errorf("something wrong happens")
}

func CatchUpdateError(err error, rowsAffected int64) error {
	if err == nil {
		if rowsAffected == 0 {
			return &NotFoundError{}
		}
		return nil
	}
	log.Println("[ERROR] ", err.Error())
	return fmt.Errorf("something wrong happens")
}

type NotFoundError struct{}

func (m *NotFoundError) Error() string {
	return "item not found"
}
