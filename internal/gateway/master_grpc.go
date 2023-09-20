package gateway

import (
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/master_service"
	"google.golang.org/grpc"
)

type MasterRepository struct {
	cli gen.MasterServiceV1Client
}

func (r *MasterRepository) CheckIfMasterExists(ctx context.Context, masterID int64) error {
	_, err := r.cli.GetMasterByID(ctx, &gen.GetMasterByIDRequestV1{Id: masterID})
	if err != nil {
		return err
	}

	return nil
}

var _ MasterGateway = (*MasterRepository)(nil)

func NewMasterRepository(conn *grpc.ClientConn) *MasterRepository {
	cli := gen.NewMasterServiceV1Client(conn)
	return &MasterRepository{cli: cli}
}
