package gateway

import (
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/admin_service"
	"google.golang.org/grpc"
)

type AdminRepository struct {
	cli gen.AdminServiceV1Client
}

func (r *AdminRepository) CheckIfAdminExists(ctx context.Context, id string) error {
	_, err := r.cli.GetByID(ctx, &gen.GetByIDV1{ID: id})
	if err != nil {
		return err
	}

	return nil
}

var _ AdminGateway = (*AdminRepository)(nil)

func NewAdminRepository(conn *grpc.ClientConn) *AdminRepository {
	cli := gen.NewAdminServiceV1Client(conn)
	return &AdminRepository{cli: cli}
}
