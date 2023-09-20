package gateway

import "context"

type (
	Registry interface {
		EntityManager
		WithTx(context.Context, func(EntityManager) error) error
	}
)
