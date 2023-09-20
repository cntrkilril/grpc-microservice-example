package entity

type (
	TimeCell struct {
		ID            int64  `db:"id"`
		MasterID      int64  `db:"master_id" validate:"required,gte=1"`
		Date          string `db:"date" validate:"required,gte=1"`
		StartTime     string `db:"start_time" validate:"required,gte=1"`
		EndTime       string `db:"end_time" validate:"required,gte=1"`
		IsFree        *bool  `db:"is_free" validate:"required"`
		AppointmentID *int   `db:"appointment_id"`
	}

	GetAvailableTimeByServiceIDAndDateDTO struct {
		ServiceID int64  `db:"service_id" validate:"required,gte=1"`
		Date      string `db:"date" validate:"required,gte=1"`
	}

	AvailableTimeArray struct {
		TimeArray []AvailableTime
	}

	AvailableTime struct {
		StartTime string
		EndTime   string
	}
)
