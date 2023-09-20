package gateway

import (
	"appointment-service/internal/entity"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
)

type AppointmentRepository struct {
	db queryRunner
}

func (r *AppointmentRepository) Save(ctx context.Context, params entity.Appointment) (result entity.Appointment, err error) {
	q := `
		INSERT INTO
		    appointment (client_id, master_id, service_id, start_time, end_time, date, is_confirmed, cancelled_at, cancel_reason, cancelled_by, cancelled_by_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, master_id, client_id, service_id, start_time, end_time, date, is_confirmed, cancelled_at, cancel_reason, cancelled_by, cancelled_by_id;
		`

	err = r.db.GetContext(ctx, &result, q, params.ClientID, params.MasterID, params.ServiceID, params.StartTime, params.EndTime,
		params.Date, params.IsConfirmed, &params.CancelledAt, &params.CancelReason, &params.CancelledBy,
		&params.CancelledByID)
	if err != nil {
		return entity.Appointment{}, err
	}
	return result, nil
}

func (r *AppointmentRepository) Cancel(ctx context.Context, params entity.CancelAppointmentDTO) (result entity.Appointment, err error) {
	q := `
		UPDATE appointment
		SET
		    cancelled_at=now(),
		    cancel_reason=$2,
		    cancelled_by=$3,
		    cancelled_by_id=$4
		WHERE id=$1
		RETURNING id, client_id, master_id, service_id, start_time, end_time, date, is_confirmed, cancelled_at, cancel_reason, cancelled_by, cancelled_by_id;
		`

	err = r.db.GetContext(ctx, &result, q, params.ID, params.CancelReason, params.CancelledBy, params.CancelledByID)
	if err != nil {
		return entity.Appointment{}, err
	}
	return result, nil
}

func (r *AppointmentRepository) Confirm(ctx context.Context, params entity.ConfirmAppointmentDTO) (result entity.Appointment, err error) {
	q := `
		UPDATE appointment
		SET
		    is_confirmed=true
		WHERE id=$1 and client_id=$2
		RETURNING id, client_id, master_id, service_id, start_time, end_time, date, is_confirmed, cancelled_at, cancel_reason, cancelled_by, cancelled_by_id;
		`

	err = r.db.GetContext(ctx, &result, q, params.ID, params.ClientID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.ErrAppointmentNotFound
		}
		return entity.Appointment{}, err
	}
	return result, nil
}

func (r *AppointmentRepository) FindByMasterID(ctx context.Context, params entity.GetAppointmentByMasterIDDTO) (result []entity.Appointment, err error) {
	q := `
		SELECT id, master_id, client_id, service_id, start_time, end_time, date, is_confirmed, cancelled_at, cancel_reason, cancelled_by, cancelled_by_id
		FROM appointment
		WHERE master_id=$1
		LIMIT $2 OFFSET $3;
		`

	err = r.db.SelectContext(ctx, &result, q, params.MasterID, params.Limit, params.Offset)
	if err != nil {
		return []entity.Appointment{}, err
	}

	return result, nil
}

func (r *AppointmentRepository) FindByClientID(ctx context.Context, params entity.GetAppointmentByClientIDDTO) (result []entity.Appointment, err error) {
	q := `
		SELECT id, client_id, service_id, start_time, end_time, date, is_confirmed, cancelled_at, cancel_reason, cancelled_by, cancelled_by_id
		FROM appointment
		WHERE client_id=$1
		LIMIT $2 OFFSET $3;
		`

	err = r.db.SelectContext(ctx, &result, q, params.ClientID, params.Limit, params.Offset)
	if err != nil {
		return []entity.Appointment{}, err
	}

	return result, nil
}

func (r *AppointmentRepository) CountByMasterID(ctx context.Context, masterID int64) (int64, error) {
	q := `
		SELECT COUNT(*) AS count FROM appointment WHERE master_id=$1;
	`

	var count int64
	err := r.db.GetContext(ctx, &count, q, masterID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AppointmentRepository) CountByClientID(ctx context.Context, clientID int64) (int64, error) {
	q := `
		SELECT COUNT(*) AS count FROM appointment WHERE client_id=$1;
	`

	var count int64
	err := r.db.GetContext(ctx, &count, q, clientID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

var _ AppointmentGateway = (*AppointmentRepository)(nil)

func NewAppointmentRepository(db *sqlx.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}
