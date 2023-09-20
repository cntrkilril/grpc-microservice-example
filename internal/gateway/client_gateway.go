package gateway

import "context"

type ClientGateway interface {
	CheckIfClientExists(context.Context, int64) error
}
