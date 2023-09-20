package gateway

import (
	"appointment-service/internal/entity"
	"context"
)

type ServicesGateway interface {
	CheckIfServiceExists(context.Context, int64) (entity.Service, error)
}
