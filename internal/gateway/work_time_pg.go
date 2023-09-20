package gateway

import (
	"appointment-service/internal/entity"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
)

type (
	WorkTimeRepository struct {
		db queryRunner
	}
)

func (r *WorkTimeRepository) Save(ctx context.Context, params entity.WorkTime) (result entity.WorkTime, err error) {
	q := `
		INSERT INTO
		    work_time (master_id, start_time, end_time, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, master_id, start_time, end_time, date;
		`

	err = r.db.GetContext(ctx, &result, q, params.MasterID, params.StartTime, params.EndTime, params.Date)
	if err != nil {
		return entity.WorkTime{}, err
	}
	return result, nil
}

func (r *WorkTimeRepository) FindByMasterID(ctx context.Context, params entity.GetWorkTimeByMasterIDDTO) (result []entity.WorkTime, err error) {
	q := `
		SELECT id, master_id, start_time, end_time, date
		FROM work_time
		WHERE master_id=$1
		LIMIT $2 OFFSET $3;
		`

	err = r.db.SelectContext(ctx, &result, q, params.MasterID, params.Limit, params.Offset)
	if err != nil {
		return []entity.WorkTime{}, err
	}

	return result, nil
}

func (r *WorkTimeRepository) FindByMasterIDAndDate(ctx context.Context, params entity.GetWorkTimeByMasterIDAndDateDTO) (result entity.WorkTime, err error) {
	q := `
		SELECT id, master_id, start_time, end_time, date
		FROM work_time
		WHERE master_id=$1 and date=$2;
		`

	err = r.db.GetContext(ctx, &result, q, params.MasterID, params.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.WorkTime{}, errors.ErrWorkTimeNotFound
		}
		return entity.WorkTime{}, err
	}

	return result, nil
}

func (r *WorkTimeRepository) CountByMasterID(ctx context.Context, masterID int64) (int64, error) {
	q := `
		SELECT COUNT(*) AS count FROM work_time WHERE master_id=$1;
	`

	var count int64
	err := r.db.GetContext(ctx, &count, q, masterID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *WorkTimeRepository) Update(ctx context.Context, params entity.UpdateWorkTimeDTO) (result entity.WorkTime, err error) {
	q := `
		UPDATE work_time
		SET
		    start_time=$3,
		    end_time=$4
		WHERE master_id=$1 and date=$2
		RETURNING id, master_id, start_time, end_time, date;
		`

	err = r.db.GetContext(ctx, &result, q, params.MasterID, params.Date, params.StartTime, params.EndTime)
	if err != nil {
		return entity.WorkTime{}, err
	}
	return result, nil
}

func (r *WorkTimeRepository) DeleteByMasterIDAndDate(ctx context.Context, params entity.DeleteWorkTimeDTO) (err error) {
	q := `
		DELETE
		FROM work_time
		WHERE master_id=$1 and date=$2;
		`
	_, err = r.db.ExecContext(ctx, q, params.MasterID, params.Date)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

var _ WorkTimeGateway = (*WorkTimeRepository)(nil)

func NewWorkTimeRepository(db *sqlx.DB) *WorkTimeRepository {
	return &WorkTimeRepository{db: db}
}
