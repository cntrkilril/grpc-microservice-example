package entity

import "google.golang.org/genproto/googleapis/type/decimal"

type (
	Service struct {
		ID          int64
		Name        string
		Description *string
		Price       decimal.Decimal
		Duration    *string
		CategoryID  int64
		MasterID    int64
	}
)
