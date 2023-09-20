package service

import (
	"appointment-service/internal/entity"
	"appointment-service/internal/gateway"
	"context"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
	"strconv"
	"strings"
)

type (
	UpdateService struct {
		clientRepo     gateway.ClientGateway
		masterRepo     gateway.MasterGateway
		repos          gateway.Registry
		servicesRepo   gateway.ServicesGateway
		adminRepo      gateway.AdminGateway
		spacerTimeCell int
	}

	UpdateInteractor interface {
		CancelAppointment(context.Context, entity.CancelAppointmentDTO) (entity.Appointment, error)
		ConfirmAppointment(context.Context, entity.ConfirmAppointmentDTO) (entity.Appointment, error)
		UpdateWorkTime(context.Context, entity.UpdateWorkTimeDTO) (entity.WorkTime, error)
	}
)

func (s *UpdateService) CancelAppointment(ctx context.Context, p entity.CancelAppointmentDTO) (entity.Appointment, error) {

	var err error

	switch p.CancelledBy {
	case "admin":
		err = s.adminRepo.CheckIfAdminExists(ctx, strconv.Itoa(int(p.CancelledByID)))
	case "master":
		err = s.masterRepo.CheckIfMasterExists(ctx, p.CancelledByID)
	case "client":
		err = s.clientRepo.CheckIfClientExists(ctx, p.CancelledByID)
	default:
		return entity.Appointment{}, errors.HandleServiceError(errors.ErrValidationError)
	}

	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	var result entity.Appointment

	err = s.repos.WithTx(ctx, func(m gateway.EntityManager) (err error) {
		result, err = m.Appointment().Cancel(ctx, p)

		if err != nil {
			return err
		}

		_, err = m.TimeCell().CancelByAppointmentID(ctx, p.ID)

		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	return result, nil

}

func (s *UpdateService) ConfirmAppointment(ctx context.Context, p entity.ConfirmAppointmentDTO) (entity.Appointment, error) {

	var result entity.Appointment

	var err error

	result, err = s.repos.Appointment().Confirm(ctx, p)

	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	return result, nil

}

func (s *UpdateService) UpdateWorkTime(ctx context.Context, p entity.UpdateWorkTimeDTO) (entity.WorkTime, error) {

	timeCells, err := s.repos.TimeCell().FindByDateAndMasterID(ctx, gateway.FindTimeCellsByDateAnsMasterIDDTO{MasterID: p.MasterID, Date: p.Date})
	if err != nil {
		return entity.WorkTime{}, errors.HandleServiceError(err)
	}

	updateFlag := true
	for _, timeCell := range timeCells {
		timeCellIsFree := timeCell.IsFree
		if !*timeCellIsFree {
			updateFlag = false
			break
		}
	}

	if !updateFlag {
		return entity.WorkTime{}, errors.HandleServiceError(err)
	}

	var result entity.WorkTime

	err = s.repos.WithTx(ctx, func(m gateway.EntityManager) (err error) {
		result, err = m.WorkTime().Update(ctx, p)

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

		startTimeSplit := strings.Split(p.StartTime, ":")
		startTimeHour, _ := strconv.Atoi(startTimeSplit[0])
		startTimeMinute, _ := strconv.Atoi(startTimeSplit[1])
		startTimeSum := startTimeHour*60 + startTimeMinute

		endTimeSplit := strings.Split(p.EndTime, ":")
		endTimeHour, _ := strconv.Atoi(endTimeSplit[0])
		endTimeMinute, _ := strconv.Atoi(endTimeSplit[1])
		endTimeSum := endTimeHour*60 + endTimeMinute

		for startTimeSum < endTimeSum {
			newEndTimeSum := startTimeSum + s.spacerTimeCell
			_, err = m.TimeCell().Save(ctx, entity.TimeCell{
				MasterID:  p.MasterID,
				Date:      p.Date,
				StartTime: strconv.Itoa(startTimeSum/60) + ":" + strconv.Itoa(startTimeSum-startTimeSum/60*60) + ":00",
				EndTime:   strconv.Itoa(newEndTimeSum/60) + ":" + strconv.Itoa(newEndTimeSum-newEndTimeSum/60*60) + ":00",
			})
			if err != nil {
				return err
			}
			startTimeSum = startTimeSum + s.spacerTimeCell
		}

		return err
	})

	if err != nil {
		return entity.WorkTime{}, errors.HandleServiceError(err)
	}

	return result, nil
}

var _ UpdateInteractor = (*UpdateService)(nil)

func NewUpdateService(repos gateway.Registry, clientRepo gateway.ClientGateway, masterRepo gateway.MasterGateway, adminRepo gateway.AdminGateway, servicesRepo gateway.ServicesGateway, spacerTimeCell int) *UpdateService {
	return &UpdateService{
		clientRepo:     clientRepo,
		masterRepo:     masterRepo,
		adminRepo:      adminRepo,
		servicesRepo:   servicesRepo,
		spacerTimeCell: spacerTimeCell,
		repos:          repos,
	}
}
