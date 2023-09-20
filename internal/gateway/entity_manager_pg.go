package gateway

type (
	pgManager struct {
		appointment *AppointmentRepository
		workTime    *WorkTimeRepository
		timeCell    *TimeCellRepository
	}
)

func (m *pgManager) Appointment() AppointmentGateway {
	return m.appointment
}

func (m *pgManager) WorkTime() WorkTimeGateway {
	return m.workTime
}

func (m *pgManager) TimeCell() TimeCellGateway {
	return m.timeCell
}

var _ EntityManager = &pgManager{}
