package resto

import "github.com/wildangbudhi/yum-service/domain"

type RestoSchedules struct {
	ID        *domain.UUID `json:"id"`
	RestoID   *domain.UUID `json:"-"`
	DayOfWeek *int         `json:"day_of_week"`
	StartTime *int         `json:"start_time"`
	EndTime   *int         `json:"end_time"`
}

type RestoSchedulesResporitory interface {
	FetchRestoScheduleByRestoID(restoID *domain.UUID) ([]RestoSchedules, error, domain.RepositoryErrorType)
	CreateOrUpdateSchedules(schedules []RestoSchedules) (error, domain.RepositoryErrorType)
}
