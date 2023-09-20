package gateway

import (
	"appointment-service/internal/entity"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
)

type (
	TimeCellRepository struct {
		db queryRunner
	}
)

func (r *TimeCellRepository) Save(ctx context.Context, params entity.TimeCell) (result entity.TimeCell, err error) {
	q := `
		INSERT INTO
		    time_cell (master_id, date, start_time, end_time, is_free)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id, master_id, date, start_time, end_time, is_free, appointment_id;
		`

	err = r.db.GetContext(ctx, &result, q, params.MasterID, params.Date, params.StartTime, params.EndTime)
	if err != nil {
		return entity.TimeCell{}, err
	}
	return result, nil
}

func (r *TimeCellRepository) FindByDateAndMasterID(ctx context.Context, params FindTimeCellsByDateAnsMasterIDDTO) (result []entity.TimeCell, err error) {
	q := `
		SELECT id, master_id, date, start_time, end_time, is_free, appointment_id
		FROM time_cell
		WHERE master_id=$1 and date=$2;
		`

	err = r.db.SelectContext(ctx, &result, q, params.MasterID, params.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.ErrTimeCellNotFound
		}
		return []entity.TimeCell{}, err
	}

	return result, nil
}

func (r *TimeCellRepository) FindByDateAndMasterIDSortByStartTime(ctx context.Context, params FindTimeCellsByDateAnsMasterIDDTO) (result []entity.TimeCell, err error) {
	q := `
		SELECT id, master_id, date, start_time, end_time, is_free, appointment_id
		FROM time_cell
		WHERE master_id=$1 and date=$2
		ORDER BY start_time;
		`

	err = r.db.SelectContext(ctx, &result, q, params.MasterID, params.Date)
	if err != nil {
		return []entity.TimeCell{}, err
	}

	return result, nil
}

func (r *TimeCellRepository) UpdateFreeByID(ctx context.Context, params UpdateTimeCellIsFreeByIDDTO) (result entity.TimeCell, err error) {
	q := `
		UPDATE time_cell
		SET is_free=$2, appointment_id=$3
		WHERE id=$1
		RETURNING id, master_id, date, start_time, end_time, is_free, appointment_id;
		`

	err = r.db.GetContext(ctx, &result, q, params.ID, params.IsFree, params.AppointmentID)
	if err != nil {
		return entity.TimeCell{}, err
	}

	return result, nil
}

func (r *TimeCellRepository) CancelByAppointmentID(ctx context.Context, appointmentID int64) (result []entity.TimeCell, err error) {
	q := `
		UPDATE time_cell
		SET is_free=True, appointment_id=null
		WHERE appointment_id=$1
		RETURNING id, master_id, date, start_time, end_time, is_free, appointment_id;
		`

	err = r.db.SelectContext(ctx, &result, q, appointmentID)
	if err != nil {
		return []entity.TimeCell{}, err
	}

	return result, nil
}

func (r *TimeCellRepository) DeleteByMasterIDAndDate(ctx context.Context, params DeleteTimeCellsByMasterIDAndDateDTO) (err error) {
	q := `
		DELETE
		FROM time_cell
		WHERE master_id=$1 and date=$2;
		`
	_, err = r.db.ExecContext(ctx, q, params.MasterID, params.Date)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

var _ TimeCellGateway = (*TimeCellRepository)(nil)

func NewTimeCellRepository(db *sqlx.DB) *TimeCellRepository {
	return &TimeCellRepository{db: db}
}
