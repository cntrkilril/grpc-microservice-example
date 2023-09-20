package gateway

import (
	"context"
)

type MasterGateway interface {
	CheckIfMasterExists(context.Context, int64) error
}
