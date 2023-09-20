package entity

type (
	WorkTime struct {
		ID        int64  `db:"id"`
		MasterID  int64  `db:"master_id" validate:"required,gte=1"`
		Date      string `db:"date" validate:"required,gte=1"`
		StartTime string `db:"start_time" validate:"required,gte=1"`
		EndTime   string `db:"end_time" validate:"required,gte=1"`
	}

	WorkTimeArray struct {
		Count     int64
		WorkTimes []WorkTime
	}

	GetWorkTimeByMasterIDDTO struct {
		MasterID int64 `db:"master_id" validate:"required,gte=1"`
		PaginationRequest
	}

	GetWorkTimeByMasterIDAndDateDTO struct {
		MasterID int64  `db:"master_id" validate:"required,gte=1"`
		Date     string `db:"date" validate:"required,gte=1"`
	}

	UpdateWorkTimeDTO struct {
		MasterID  int64  `db:"master_id" validate:"required,gte=1"`
		Date      string `db:"date" validate:"required,gte=1"`
		StartTime string `db:"start_time" validate:"required,gte=1"`
		EndTime   string `db:"end_time" validate:"required,gte=1"`
	}

	DeleteWorkTimeDTO struct {
		MasterID int64  `db:"master_id" validate:"required,gte=1"`
		Date     string `db:"date" validate:"required,gte=1"`
	}
)
