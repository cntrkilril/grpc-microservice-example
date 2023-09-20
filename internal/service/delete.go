package service

import (
	"appointment-service/internal/entity"
	"appointment-service/internal/gateway"
	"context"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
)

type (
	DeleteService struct {
		repos gateway.Registry
	}

	DeleteInteractor interface {
		DeleteWorkTimeByMasterIDAndDate(context.Context, entity.DeleteWorkTimeDTO) error
	}
)

func (s *DeleteService) DeleteWorkTimeByMasterIDAndDate(ctx context.Context, p entity.DeleteWorkTimeDTO) error {

	timeCells, err := s.repos.TimeCell().FindByDateAndMasterID(ctx, gateway.FindTimeCellsByDateAnsMasterIDDTO{MasterID: p.MasterID, Date: p.Date})
	if err != nil {
		return errors.HandleServiceError(err)
	}

	deleteFlag := true
	for _, timeCell := range timeCells {
		timeCellIsFree := timeCell.IsFree
		if !*timeCellIsFree {
			deleteFlag = false
			break
		}
	}

	if !deleteFlag {
		return errors.HandleServiceError(err)
	}

	err = s.repos.WithTx(ctx, func(m gateway.EntityManager) (err error) {
		err = m.WorkTime().DeleteByMasterIDAndDate(ctx, p)

		if err != nil {
			return err
		}

		err = m.TimeCell().DeleteByMasterIDAndDate(ctx, gateway.DeleteTimeCellsByMasterIDAndDateDTO{
			MasterID: p.MasterID,
			Date:     p.Date,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.HandleServiceError(err)
	}

	return nil
}

var _ DeleteInteractor = (*DeleteService)(nil)

func NewDeleteService(repos gateway.Registry) *DeleteService {
	return &DeleteService{
		repos: repos,
	}
}
