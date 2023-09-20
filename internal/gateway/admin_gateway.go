package gateway

import "context"

type AdminGateway interface {
	CheckIfAdminExists(context.Context, string) error
}
