package gateway

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type (
	PGRegistry struct {
		db *sqlx.DB
		m  *pgManager
	}
)

func (r *PGRegistry) Appointment() AppointmentGateway {
	return r.m.appointment
}

func (r *PGRegistry) WorkTime() WorkTimeGateway {
	return r.m.workTime
}

func (r *PGRegistry) TimeCell() TimeCellGateway {
	return r.m.timeCell
}

func (r *PGRegistry) WithTx(ctx context.Context, f func(EntityManager) error) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = f(&pgManager{
		appointment: &AppointmentRepository{tx},
		workTime:    &WorkTimeRepository{tx},
		timeCell:    &TimeCellRepository{tx},
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

var _ Registry = &PGRegistry{}

func NewPGRegistry(db *sqlx.DB) *PGRegistry {
	return &PGRegistry{
		db: db,
		m: &pgManager{
			appointment: &AppointmentRepository{db},
			workTime:    &WorkTimeRepository{db},
			timeCell:    &TimeCellRepository{db: db},
		},
	}
}
