package gateway

type EntityManager interface {
	Appointment() AppointmentGateway
	TimeCell() TimeCellGateway
	WorkTime() WorkTimeGateway
}
