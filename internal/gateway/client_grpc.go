package gateway

import (
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/client_service"
	"google.golang.org/grpc"
)

type ClientRepository struct {
	cli gen.ClientServiceV1Client
}

func (r *ClientRepository) CheckIfClientExists(ctx context.Context, clientID int64) error {
	_, err := r.cli.GetClientByID(ctx, &gen.GetClientByIDRequestV1{Id: clientID})
	if err != nil {
		return err
	}

	return nil
}

var _ ClientGateway = (*ClientRepository)(nil)

func NewClientRepository(conn *grpc.ClientConn) *ClientRepository {
	cli := gen.NewClientServiceV1Client(conn)
	return &ClientRepository{cli: cli}
}
