package gateway

import (
	"appointment-service/internal/entity"
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/services_service"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/grpc"
)

type ServicesRepository struct {
	cli gen.ServicesServiceV1Client
}

func (r *ServicesRepository) CheckIfServiceExists(ctx context.Context, servicesID int64) (entity.Service, error) {
	result, err := r.cli.GetServiceByID(ctx, &gen.GetServiceByIDRequestV1{ID: servicesID})
	if err != nil {
		return entity.Service{}, err
	}

	return entity.Service{
		ID:          result.ID,
		Name:        result.Name,
		Description: &result.Description,
		Price:       decimal.Decimal{Value: result.Price},
		Duration:    &result.Duration,
		CategoryID:  result.CategoryID,
		MasterID:    result.MasterID,
	}, nil
}

var _ ServicesGateway = (*ServicesRepository)(nil)

func NewServicesRepository(conn *grpc.ClientConn) *ServicesRepository {
	cli := gen.NewServicesServiceV1Client(conn)
	return &ServicesRepository{cli: cli}
}
