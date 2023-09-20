package gateway

import (
	"appointment-service/internal/entity"
	"context"
)

type (
	WorkTimeGateway interface {
		Save(context.Context, entity.WorkTime) (entity.WorkTime, error)
		FindByMasterID(context.Context, entity.GetWorkTimeByMasterIDDTO) ([]entity.WorkTime, error)
		FindByMasterIDAndDate(context.Context, entity.GetWorkTimeByMasterIDAndDateDTO) (entity.WorkTime, error)
		CountByMasterID(context.Context, int64) (int64, error)
		Update(context.Context, entity.UpdateWorkTimeDTO) (entity.WorkTime, error)
		DeleteByMasterIDAndDate(context.Context, entity.DeleteWorkTimeDTO) error
	}
)
