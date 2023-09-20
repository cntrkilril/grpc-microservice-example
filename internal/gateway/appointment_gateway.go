package gateway

import (
	"appointment-service/internal/entity"
	"context"
)

type AppointmentGateway interface {
	Save(context.Context, entity.Appointment) (entity.Appointment, error)
	Cancel(context.Context, entity.CancelAppointmentDTO) (entity.Appointment, error)
	Confirm(context.Context, entity.ConfirmAppointmentDTO) (entity.Appointment, error)
	FindByMasterID(context.Context, entity.GetAppointmentByMasterIDDTO) ([]entity.Appointment, error)
	CountByMasterID(context.Context, int64) (int64, error)
	FindByClientID(context.Context, entity.GetAppointmentByClientIDDTO) ([]entity.Appointment, error)
	CountByClientID(context.Context, int64) (int64, error)
}
