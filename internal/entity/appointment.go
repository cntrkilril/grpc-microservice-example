package entity

type (
	Appointment struct {
		ID            int64   `db:"id"`
		MasterID      int64   `db:"master_id" validate:"required,gte=1"`
		ClientID      int64   `db:"client_id" validate:"required,gte=1"`
		ServiceID     int64   `db:"service_id" validate:"required,gte=1"`
		StartTime     string  `db:"start_time" validate:"required,gte=1"`
		EndTime       string  `db:"end_time" validate:"required,gte=1"`
		Date          string  `db:"date" validate:"required,gte=1"`
		IsConfirmed   bool    `db:"is_confirmed" validate:"required"`
		CancelledAt   *string `db:"cancelled_at"`
		CancelReason  *string `db:"cancel_reason"`
		CancelledBy   *string `db:"cancelled_by"`
		CancelledByID *int64  `db:"cancelled_by_id"`
	}

	AppointmentArray struct {
		Count        int64
		Appointments []Appointment
	}

	CancelAppointmentDTO struct {
		ID            int64  `db:"id" validate:"required"`
		CancelReason  string `db:"cancel_reason" validate:"required"`
		CancelledBy   string `db:"cancelled_by" validate:"required"`
		CancelledByID int64  `db:"cancelled_by_id" validate:"required"`
	}

	ConfirmAppointmentDTO struct {
		ID       int64 `db:"id" validate:"required"`
		ClientID int64 `db:"client_id" validate:"required"`
	}

	GetAppointmentByMasterIDDTO struct {
		MasterID int64 `db:"master_id" validate:"required,gte=1"`
		PaginationRequest
	}

	GetAppointmentByClientIDDTO struct {
		ClientID int64 `db:"client_id" validate:"required,gte=1"`
		PaginationRequest
	}

	PaginationRequest struct {
		Limit  int64 `validate:"gte=1"`
		Offset int64 `validate:"gte=0"`
	}
)
