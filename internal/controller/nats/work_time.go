package nats

import (
	"appointment-service/internal/entity"
	"appointment-service/internal/service"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type WorkTimeController struct {
	js            nats.JetStreamContext
	createService service.CreateInteractor
	lgr           *zap.SugaredLogger
}

const (
	StreamName                   = "APPOINTMENT"
	CreateWorkTimeSubjectSubject = "create_work_time"
)

func (c *WorkTimeController) createWorkTime(m *nats.Msg) {
	var data entity.WorkTime
	err := json.Unmarshal(m.Data, &data)
	if err != nil {
		c.lgr.Error(err)
		return
	}

	_, err = c.createService.CreateWorkTime(context.Background(), data)
	if err != nil {
		c.lgr.Error(err)
		return
	}

	err = m.AckSync()
	if err != nil {
		c.lgr.Error(err)
		return
	}
}

func (c *WorkTimeController) Subscribe() error {
	_, err := c.js.Subscribe(
		CreateWorkTimeSubjectSubject,
		c.createWorkTime,
		nats.Durable(CreateWorkTimeSubjectSubject),
		nats.BindStream(StreamName),
	)
	if err != nil {
		return err
	}
	return nil
}

func NewWorkTimeController(
	js nats.JetStreamContext,
	createService service.CreateInteractor,
	lgr *zap.SugaredLogger,
) *WorkTimeController {
	return &WorkTimeController{
		js:            js,
		createService: createService,
		lgr:           lgr,
	}
}
