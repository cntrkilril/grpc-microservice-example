package gateway

import (
	"appointment-service/internal/entity"
	"context"
)

type (
	FindTimeCellsByDateAnsMasterIDDTO struct {
		MasterID int64
		Date     string
	}

	UpdateTimeCellIsFreeByIDDTO struct {
		IsFree        bool
		ID            int64
		AppointmentID int64
	}

	DeleteTimeCellsByMasterIDAndDateDTO struct {
		MasterID int64
		Date     string
	}

	TimeCellGateway interface {
		Save(context.Context, entity.TimeCell) (entity.TimeCell, error)
		FindByDateAndMasterID(context.Context, FindTimeCellsByDateAnsMasterIDDTO) ([]entity.TimeCell, error)
		FindByDateAndMasterIDSortByStartTime(context.Context, FindTimeCellsByDateAnsMasterIDDTO) ([]entity.TimeCell, error)
		UpdateFreeByID(context.Context, UpdateTimeCellIsFreeByIDDTO) (entity.TimeCell, error)
		CancelByAppointmentID(context.Context, int64) ([]entity.TimeCell, error)
		DeleteByMasterIDAndDate(context.Context, DeleteTimeCellsByMasterIDAndDateDTO) error
	}
)
