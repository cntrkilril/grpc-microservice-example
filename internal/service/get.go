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
	GetService struct {
		appointmentRepo gateway.AppointmentGateway
		workTimeRepo    gateway.WorkTimeGateway
		servicesRepo    gateway.ServicesGateway
		timeCellRepo    gateway.TimeCellGateway
		spacerTimeCell  int
	}

	GetInteractor interface {
		GetAppointmentByClientID(context.Context, entity.GetAppointmentByClientIDDTO) (entity.AppointmentArray, error)
		GetAppointmentByMasterID(context.Context, entity.GetAppointmentByMasterIDDTO) (entity.AppointmentArray, error)
		GetWorkTimeByMasterID(context.Context, entity.GetWorkTimeByMasterIDDTO) (entity.WorkTimeArray, error)
		GetAvailableTimeByServiceIDAndDate(context.Context, entity.GetAvailableTimeByServiceIDAndDateDTO) (entity.AvailableTimeArray, error)
	}
)

func (s *GetService) GetAppointmentByClientID(ctx context.Context, p entity.GetAppointmentByClientIDDTO) (entity.AppointmentArray, error) {

	result, err := s.appointmentRepo.FindByClientID(ctx, p)
	if err != nil {
		return entity.AppointmentArray{}, errors.HandleServiceError(err)
	}

	count, err := s.appointmentRepo.CountByClientID(ctx, p.ClientID)
	if err != nil {
		return entity.AppointmentArray{}, errors.HandleServiceError(err)
	}

	return entity.AppointmentArray{Count: count, Appointments: result}, nil
}

func (s *GetService) GetAppointmentByMasterID(ctx context.Context, p entity.GetAppointmentByMasterIDDTO) (entity.AppointmentArray, error) {

	result, err := s.appointmentRepo.FindByMasterID(ctx, p)
	if err != nil {
		return entity.AppointmentArray{}, errors.HandleServiceError(err)
	}

	count, err := s.appointmentRepo.CountByMasterID(ctx, p.MasterID)
	if err != nil {
		return entity.AppointmentArray{}, errors.HandleServiceError(err)
	}

	return entity.AppointmentArray{Count: count, Appointments: result}, nil
}

func (s *GetService) GetWorkTimeByMasterID(ctx context.Context, p entity.GetWorkTimeByMasterIDDTO) (entity.WorkTimeArray, error) {

	result, err := s.workTimeRepo.FindByMasterID(ctx, p)
	if err != nil {
		return entity.WorkTimeArray{}, errors.HandleServiceError(err)
	}

	count, err := s.workTimeRepo.CountByMasterID(ctx, p.MasterID)
	if err != nil {
		return entity.WorkTimeArray{}, errors.HandleServiceError(err)
	}

	return entity.WorkTimeArray{Count: count, WorkTimes: result}, nil
}

func (s *GetService) GetAvailableTimeByServiceIDAndDate(ctx context.Context, p entity.GetAvailableTimeByServiceIDAndDateDTO) (entity.AvailableTimeArray, error) {

	serviceRes, err := s.servicesRepo.CheckIfServiceExists(ctx, p.ServiceID)
	if err != nil {
		return entity.AvailableTimeArray{}, errors.HandleServiceError(err)
	}

	timeCells, err := s.timeCellRepo.FindByDateAndMasterIDSortByStartTime(ctx, gateway.FindTimeCellsByDateAnsMasterIDDTO{MasterID: serviceRes.MasterID, Date: p.Date})
	//timeCells, err := s.timeCellRepo.FindByDateAndMasterIDSortByStartTime(ctx, gateway.FindTimeCellsByDateAnsMasterIDDTO{MasterID: 1, Date: p.Date})
	if err != nil {
		return entity.AvailableTimeArray{}, errors.HandleServiceError(err)
	}

	durationSplit := strings.Split(*serviceRes.Duration, ":")
	//durationSplit := strings.Split("01:30:00", ":")
	hourDuration, _ := strconv.Atoi(durationSplit[0])
	minuteDuration, _ := strconv.Atoi(durationSplit[1])
	needCountTimeCell := (hourDuration*60 + minuteDuration) / s.spacerTimeCell

	if len(timeCells) == 0 {
		return entity.AvailableTimeArray{
			TimeArray: []entity.AvailableTime{},
		}, nil
	}

	currentCountTimeCell := 0
	var firstTimeCellInChain entity.TimeCell
	var result entity.AvailableTimeArray
	for _, timeCell := range timeCells {
		timeCellIsFree := timeCell.IsFree
		if *timeCellIsFree {
			if currentCountTimeCell == 0 {
				firstTimeCellInChain = timeCell
			}
			currentCountTimeCell++
			if currentCountTimeCell == needCountTimeCell {
				result.TimeArray = append(result.TimeArray, entity.AvailableTime{StartTime: firstTimeCellInChain.StartTime, EndTime: timeCell.EndTime})
				currentCountTimeCell = 0
			}
		} else {
			currentCountTimeCell = 0
		}
	}

	return result, nil
}

var _ GetInteractor = (*GetService)(nil)

func NewGetService(appointmentRepo gateway.AppointmentGateway, workTimeRepo gateway.WorkTimeGateway, servicesRepo gateway.ServicesGateway, timeCellRepo gateway.TimeCellGateway, spacerTimeCell int) *GetService {
	return &GetService{
		appointmentRepo: appointmentRepo,
		workTimeRepo:    workTimeRepo,
		servicesRepo:    servicesRepo,
		timeCellRepo:    timeCellRepo,
		spacerTimeCell:  spacerTimeCell,
	}
}
