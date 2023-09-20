package v1

import (
	"appointment-service/internal/entity"
	"appointment-service/internal/service"
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/appointment_service"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
	"gitlab.com/d1zero-online-booking/common/pkg/govalidator"
)

type TimeCellController struct {
	gen.UnimplementedTimeCellServiceV1Server
	getService service.GetInteractor
	val        *govalidator.Validator
}

func (c *TimeCellController) GetAvailableTimeByServiceIDAndDate(ctx context.Context,
	req *gen.GetAvailableTimeByServiceIDAndDateRequestV1) (*gen.GetAvailableTimeByServiceIDAndDateResponseV1, error) {

	params := entity.GetAvailableTimeByServiceIDAndDateDTO{
		ServiceID: req.GetServiceId(),
		Date:      req.GetDate(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.GetAvailableTimeByServiceIDAndDateResponseV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	res, err := c.getService.GetAvailableTimeByServiceIDAndDate(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := gen.GetAvailableTimeByServiceIDAndDateResponseV1{}
	for _, availableTime := range res.TimeArray {
		result.TimeArray = append(result.TimeArray, &gen.AvailableTimeEntityV1{
			StartTime: availableTime.StartTime,
			EndTime:   availableTime.EndTime,
		})
	}

	return &result, nil
}

func NewTimeCellController(
	getService service.GetInteractor,
	val *govalidator.Validator,
) *TimeCellController {
	return &TimeCellController{
		getService: getService,
		val:        val,
	}
}
