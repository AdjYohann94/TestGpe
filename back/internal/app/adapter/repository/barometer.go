package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/domain"
)

type BarometerRepository struct {
	DB *gorm.DB
}

const statsYearPerWorkCategory = `select avg(barometers.score) as score, count(*) as entries, type, date_part('year', date) as year, work_category_id from barometers inner join users u on u.id = barometers.user_id inner join work_categories wc on wc.id = u.work_category_id group by type, year, work_category_id`
const statsMonthPerWorkCategory = `select avg(barometers.score) as score, count(*) as entries, type, date_part('year', date) as year, date_part('month', date) as month, work_category_id from barometers inner join users u on u.id = barometers.user_id inner join work_categories wc on wc.id = u.work_category_id group by type, year, month, work_category_id`
const statsDayPerWorkCategory = `select avg(barometers.score) as score, count(*) as entries, type,date_part('year', date) as year, date_part('month', date) as month, date_part('day', date) as day, work_category_id from barometers inner join users u on u.id = barometers.user_id inner join work_categories wc on wc.id = u.work_category_id group by type, year, month, day, work_category_id`

const statsYear = `select avg(barometers.score) as score, count(*) as entries, type, date_part('year', date) as year from barometers group by type, year`
const statsMonth = `select avg(barometers.score) as score, count(*) as entries, type, date_part('year', date) as year, date_part('month', date) as month from barometers group by type, year, month`
const statsDay = `select avg(barometers.score) as score, count(*) as entries, type, date_part('year', date) as year, date_part('month', date) as month, date_part('day', date) as day from barometers group by type, year, month, day`

// DefineForDay the barometer will save or update the barometer daily value for the user and type.
// If the barometer is already defined for the current day, only score will be updated.
// The check takes the current time because it was a daily barometer.
func (mr BarometerRepository) DefineForDay(barometer *model.Barometer) (*model.Barometer, error) {
	barometer.Date = getCurrentDate()
	result := mr.DB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}, {Name: "type"}},
			DoUpdates: clause.AssignmentColumns([]string{"score"}),
		},
	).
		Create(&barometer)
	err := CatchCreateError(result.Error, result.RowsAffected)
	if err != nil {
		return nil, err
	}

	return barometer, nil
}

// RetrieveForDayUser will return all barometers saved for the user today.
func (mr BarometerRepository) RetrieveForDayUser(userID uint) ([]*model.Barometer, error) {
	var barometers []*model.Barometer
	result := mr.DB.Where(model.Barometer{UserID: userID, Date: getCurrentDate()}).Find(&barometers)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return barometers, nil
}

// GetStats will return all barometer statistics. The scopeCategories filter defines if response must be separated by user
// work categories. The scopeTime filter defines the time range of the period to compute.
// For example is scopeTime value is month, barometer scores are average on the month.
// Output stats are the average score and the number of barometer filled to generate this average score.
func (mr BarometerRepository) GetStats(scopedCategories bool, scopeTime string) ([]*domain.BarometerStat, error) {
	var query string
	if scopedCategories {
		if scopeTime == "day" {
			query = statsDayPerWorkCategory
		} else if scopeTime == "month" {
			query = statsMonthPerWorkCategory
		} else {
			query = statsYearPerWorkCategory
		}
	} else {
		if scopeTime == "day" {
			query = statsDay
		} else if scopeTime == "month" {
			query = statsMonth
		} else {
			query = statsYear
		}
	}
	var stats []*domain.BarometerStat
	result := mr.DB.Raw(query).Scan(&stats)
	err := CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func getCurrentDate() *time.Time {
	y, m, d := time.Now().Date()
	currentDay, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%d-%d", y, m, d))
	return &currentDay
}
