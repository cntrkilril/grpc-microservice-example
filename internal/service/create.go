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
	CreateService struct {
		clientRepo     gateway.ClientGateway
		servicesRepo   gateway.ServicesGateway
		masterRepo     gateway.MasterGateway
		repos          gateway.Registry
		spacerTimeCell int
	}

	CreateInteractor interface {
		CreateAppointment(context.Context, entity.Appointment) (entity.Appointment, error)
		CreateWorkTime(ctx context.Context, time entity.WorkTime) (entity.WorkTime, error)
	}
)

func (s *CreateService) CreateAppointment(ctx context.Context, p entity.Appointment) (entity.Appointment, error) {

	err := s.clientRepo.CheckIfClientExists(ctx, p.ClientID)
	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	serviceRes, err := s.servicesRepo.CheckIfServiceExists(ctx, p.ServiceID)
	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	if serviceRes.MasterID != p.MasterID {
		return entity.Appointment{}, errors.HandleServiceError(errors.ErrUnknown)
	}

	timeCells, err := s.repos.TimeCell().FindByDateAndMasterID(ctx, gateway.FindTimeCellsByDateAnsMasterIDDTO{MasterID: p.MasterID, Date: p.Date})
	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	durationSplit := strings.Split(*serviceRes.Duration, ":")
	//durationSplit := strings.Split("00:30:00", ":")
	hourDuration, _ := strconv.Atoi(durationSplit[0])
	minuteDuration, _ := strconv.Atoi(durationSplit[1])
	needCountTimeCell := (hourDuration*60 + minuteDuration) / s.spacerTimeCell

	startTimeSplit := strings.Split(p.StartTime, ":")
	startTimeHour, _ := strconv.Atoi(startTimeSplit[0])
	startTimeMinute, _ := strconv.Atoi(startTimeSplit[1])
	startTimeSum := startTimeHour*60 + startTimeMinute

	i := 0
	var needTimeCells []int
	for i < needCountTimeCell {
		needTimeCells = append(needTimeCells, startTimeSum+s.spacerTimeCell*i)
		i++
	}

	countTimeCell := 0
	availableSaveFlag := false
	var needTimeCellsIDs []int64
	for _, timeCell := range timeCells {
		timeCellIsFree := timeCell.IsFree
		if *timeCellIsFree {
			startTimeSplit = strings.Split(timeCell.StartTime, ":")
			startTimeHour, _ = strconv.Atoi(startTimeSplit[0])
			startTimeMinute, _ = strconv.Atoi(startTimeSplit[1])
			startTimeSum = startTimeHour*60 + startTimeMinute
			for _, needTimeCell := range needTimeCells {
				if needTimeCell == startTimeSum {
					countTimeCell++
					needTimeCellsIDs = append(needTimeCellsIDs, timeCell.ID)
					break
				}
			}
			if countTimeCell == needCountTimeCell {
				availableSaveFlag = true
				break
			}
		}
	}

	if !availableSaveFlag {
		return entity.Appointment{}, errors.HandleServiceError(errors.ErrTimeCellNotFound)
	}

	var result entity.Appointment

	err = s.repos.WithTx(ctx, func(m gateway.EntityManager) (err error) {
		result, err = m.Appointment().Save(ctx, p)

		if err != nil {
			return err
		}

		for _, item := range needTimeCellsIDs {
			_, err = m.TimeCell().UpdateFreeByID(ctx, gateway.UpdateTimeCellIsFreeByIDDTO{
				ID:            item,
				IsFree:        false,
				AppointmentID: result.ID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return entity.Appointment{}, errors.HandleServiceError(err)
	}

	return result, nil
}

func (s *CreateService) CreateWorkTime(ctx context.Context, p entity.WorkTime) (entity.WorkTime, error) {

	err := s.masterRepo.CheckIfMasterExists(ctx, p.MasterID)
	if err != nil {
		return entity.WorkTime{}, errors.HandleServiceError(err)
	}

	_, err = s.repos.WorkTime().FindByMasterIDAndDate(ctx, entity.GetWorkTimeByMasterIDAndDateDTO{MasterID: p.MasterID, Date: p.Date})
	if err != nil && err != errors.ErrWorkTimeNotFound {
		return entity.WorkTime{}, errors.HandleServiceError(err)
	}

	if err != errors.ErrWorkTimeNotFound {
		return entity.WorkTime{}, errors.HandleServiceError(err)
	}

	var result entity.WorkTime

	err = s.repos.WithTx(ctx, func(m gateway.EntityManager) (err error) {
		result, err = m.WorkTime().Save(ctx, p)

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

var _ CreateInteractor = (*CreateService)(nil)

func NewCreateService(repos gateway.Registry, clientRepo gateway.ClientGateway, servicesRepo gateway.ServicesGateway, masterRepo gateway.MasterGateway, spacerTimeCell int) *CreateService {
	return &CreateService{
		clientRepo:     clientRepo,
		servicesRepo:   servicesRepo,
		repos:          repos,
		masterRepo:     masterRepo,
		spacerTimeCell: spacerTimeCell,
	}
}
