package v1

import (
	"appointment-service/internal/entity"
	"appointment-service/internal/service"
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/appointment_service"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
	"gitlab.com/d1zero-online-booking/common/pkg/govalidator"
)

type AppointmentController struct {
	gen.UnimplementedAppointmentServiceV1Server
	getService    service.GetInteractor
	createService service.CreateInteractor
	updateService service.UpdateInteractor
	val           *govalidator.Validator
}

func (c *AppointmentController) CreateAppointment(ctx context.Context,
	req *gen.CreateAppointmentRequestV1) (*gen.AppointmentEntityV1, error) {

	params := entity.Appointment{
		MasterID:    req.GetMasterId(),
		ClientID:    req.GetClientId(),
		ServiceID:   req.GetServiceId(),
		StartTime:   req.GetStartTime(),
		EndTime:     req.GetEndTime(),
		Date:        req.GetDate(),
		IsConfirmed: req.GetIsConfirmed(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.AppointmentEntityV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	appointment, err := c.createService.CreateAppointment(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := &gen.AppointmentEntityV1{
		Id:          appointment.ID,
		MasterId:    appointment.MasterID,
		ClientId:    appointment.ClientID,
		ServiceId:   appointment.ServiceID,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Date:        appointment.Date,
		IsConfirmed: appointment.IsConfirmed,
	}

	return result, nil
}

func (c *AppointmentController) CancelAppointment(ctx context.Context,
	req *gen.CancelAppointmentRequestV1) (*gen.AppointmentEntityV1, error) {

	params := entity.CancelAppointmentDTO{
		ID:            req.GetId(),
		CancelReason:  req.GetCancelReason(),
		CancelledBy:   req.GetCancelledBy(),
		CancelledByID: req.GetCancelledById(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.AppointmentEntityV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	appointment, err := c.updateService.CancelAppointment(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := &gen.AppointmentEntityV1{
		Id:          appointment.ID,
		MasterId:    appointment.MasterID,
		ClientId:    appointment.ClientID,
		ServiceId:   appointment.ServiceID,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Date:        appointment.Date,
		IsConfirmed: appointment.IsConfirmed,
	}

	if appointment.CancelledAt != nil {
		result.CancelledAt = *appointment.CancelledAt
		result.CancelReason = *appointment.CancelReason
		result.CancelledBy = *appointment.CancelledBy
		result.CancelledById = *appointment.CancelledByID
	}

	return result, nil
}

func (c *AppointmentController) ConfirmAppointment(ctx context.Context,
	req *gen.ConfirmAppointmentRequestV1) (*gen.AppointmentEntityV1, error) {

	params := entity.ConfirmAppointmentDTO{
		ID:       req.GetId(),
		ClientID: req.GetClientId(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.AppointmentEntityV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	appointment, err := c.updateService.ConfirmAppointment(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := &gen.AppointmentEntityV1{
		Id:          appointment.ID,
		MasterId:    appointment.MasterID,
		ClientId:    appointment.ClientID,
		ServiceId:   appointment.ServiceID,
		StartTime:   appointment.StartTime,
		EndTime:     appointment.EndTime,
		Date:        appointment.Date,
		IsConfirmed: appointment.IsConfirmed,
	}
	if appointment.CancelledAt != nil {
		result.CancelledAt = *appointment.CancelledAt
		result.CancelReason = *appointment.CancelReason
		result.CancelledBy = *appointment.CancelledBy
		result.CancelledById = *appointment.CancelledByID
	}

	return result, nil
}

func (c *AppointmentController) GetAppointmentByMasterID(ctx context.Context,
	req *gen.GetAppointmentByMasterIDRequestV1) (*gen.GetAppointmentByMasterIDResponseV1, error) {

	params := entity.GetAppointmentByMasterIDDTO{
		MasterID: req.GetMasterId(),
		PaginationRequest: entity.PaginationRequest{
			Limit:  req.GetLimit(),
			Offset: req.GetOffset(),
		},
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.GetAppointmentByMasterIDResponseV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	res, err := c.getService.GetAppointmentByMasterID(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := gen.GetAppointmentByMasterIDResponseV1{
		Count: res.Count,
	}
	for _, appointment := range res.Appointments {

		newAppointment := &gen.AppointmentEntityV1{
			Id:          appointment.ID,
			MasterId:    appointment.MasterID,
			ClientId:    appointment.ClientID,
			ServiceId:   appointment.ServiceID,
			StartTime:   appointment.StartTime,
			EndTime:     appointment.EndTime,
			Date:        appointment.Date,
			IsConfirmed: appointment.IsConfirmed,
		}

		if appointment.CancelledAt != nil {
			newAppointment.CancelledAt = *appointment.CancelledAt
			newAppointment.CancelReason = *appointment.CancelReason
			newAppointment.CancelledBy = *appointment.CancelledBy
			newAppointment.CancelledById = *appointment.CancelledByID
		}

		result.Appointments = append(result.Appointments, newAppointment)
	}

	return &result, nil
}

func (c *AppointmentController) GetAppointmentByClientID(ctx context.Context,
	req *gen.GetAppointmentByClientIDRequestV1) (*gen.GetAppointmentByClientIDResponseV1, error) {

	params := entity.GetAppointmentByClientIDDTO{
		ClientID: req.GetClientId(),
		PaginationRequest: entity.PaginationRequest{
			Limit:  req.GetLimit(),
			Offset: req.GetOffset(),
		},
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.GetAppointmentByClientIDResponseV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	res, err := c.getService.GetAppointmentByClientID(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := gen.GetAppointmentByClientIDResponseV1{
		Count: res.Count,
	}
	for _, appointment := range res.Appointments {

		newAppointment := &gen.AppointmentEntityV1{
			Id:          appointment.ID,
			MasterId:    appointment.MasterID,
			ClientId:    appointment.ClientID,
			ServiceId:   appointment.ServiceID,
			StartTime:   appointment.StartTime,
			EndTime:     appointment.EndTime,
			Date:        appointment.Date,
			IsConfirmed: appointment.IsConfirmed,
		}

		if appointment.CancelledAt != nil {
			newAppointment.CancelledAt = *appointment.CancelledAt
			newAppointment.CancelReason = *appointment.CancelReason
			newAppointment.CancelledBy = *appointment.CancelledBy
			newAppointment.CancelledById = *appointment.CancelledByID
		}

		result.Appointments = append(result.Appointments, newAppointment)
	}

	return &result, nil
}

func NewAppointmentController(
	createService service.CreateInteractor,
	getService service.GetInteractor,
	updateService service.UpdateInteractor,
	val *govalidator.Validator,
) *AppointmentController {
	return &AppointmentController{
		createService: createService,
		getService:    getService,
		updateService: updateService,
		val:           val,
	}
}
