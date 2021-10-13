package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/resto"
)

type restoSchedulesRepository struct {
	db *sql.DB
}

func NewRestoSchedulesRepository(db *sql.DB) resto.RestoSchedulesResporitory {
	return &restoSchedulesRepository{
		db: db,
	}
}

func (repo *restoSchedulesRepository) FetchRestoScheduleByRestoID(restoID *domain.UUID) ([]resto.RestoSchedules, error, domain.RepositoryErrorType) {

	var err error

	var queryString string = `
	SELECT id, day_of_week, start_time, end_time
	FROM resto_schedules
	WHERE resto_id = ?
	`

	var queryResult *sql.Rows
	var schedules []resto.RestoSchedules = make([]resto.RestoSchedules, 0)

	queryResult, err = repo.db.Query(queryString, restoID)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	defer queryResult.Close()

	for queryResult.Next() {

		var schedule resto.RestoSchedules = resto.RestoSchedules{
			RestoID: restoID,
		}

		err = queryResult.Scan(
			&schedule.ID,
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
		)

		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Services Unavailable"), domain.RepositoryError
		}

		schedules = append(schedules, schedule)

	}

	return schedules, nil, 0

}

func (repo *restoSchedulesRepository) CreateOrUpdateSchedules(schedules []resto.RestoSchedules) (error, domain.RepositoryErrorType) {

	var err error
	var queryString string

	queryString = `
	INSERT INTO resto_schedules
	(id, resto_id, day_of_week, start_time, end_time)
	VALUES %s

	ON DUPLICATE KEY UPDATE
	day_of_week = VALUES(day_of_week),
	start_time = VALUES(start_time),
	end_time = VALUES(end_time)
	`

	var template string = ""
	var values []interface{} = make([]interface{}, 0)

	for i, row := range schedules {

		if i != 0 {
			template += ", "
		}

		template += "(?, ?, ?, ?, ?)"
		values = append(values, row.ID, row.RestoID, row.DayOfWeek, row.StartTime, row.EndTime)

	}

	queryString = fmt.Sprintf(queryString, template)

	statement, err := repo.db.Prepare(queryString)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	var res sql.Result

	res, err = statement.Exec(values...)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Services Unavailable"), domain.RepositoryError
	}

	if rowAffected == 0 {
		return fmt.Errorf("Failed to Save Partner Data"), domain.RepositoryCreateDataFailed
	}

	return nil, 0

}
