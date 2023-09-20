package v1

import (
	"appointment-service/internal/entity"
	"appointment-service/internal/service"
	"context"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/appointment_service"
	"gitlab.com/d1zero-online-booking/common/pkg/errors"
	"gitlab.com/d1zero-online-booking/common/pkg/govalidator"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkTimeController struct {
	gen.UnimplementedWorkTimeServiceV1Server
	getService    service.GetInteractor
	deleteService service.DeleteInteractor
	updateService service.UpdateInteractor
	val           *govalidator.Validator
}

func (c *WorkTimeController) GetWorkTimeByMasterID(ctx context.Context,
	req *gen.GetWorkTimeByMasterIDRequestV1) (*gen.GetWorkTimeByMasterIDResponseV1, error) {

	params := entity.GetWorkTimeByMasterIDDTO{
		MasterID: req.GetMasterId(),
		PaginationRequest: entity.PaginationRequest{
			Limit:  req.GetLimit(),
			Offset: req.GetOffset(),
		},
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.GetWorkTimeByMasterIDResponseV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	res, err := c.getService.GetWorkTimeByMasterID(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	result := gen.GetWorkTimeByMasterIDResponseV1{
		Count:     res.Count,
		WorkTimes: make([]*gen.WorkTimeEntityV1, 0, res.Count),
	}
	for _, workTime := range res.WorkTimes {
		result.WorkTimes = append(result.WorkTimes, &gen.WorkTimeEntityV1{
			Id:        workTime.ID,
			MasterId:  workTime.MasterID,
			StartTime: workTime.StartTime,
			EndTime:   workTime.EndTime,
			Date:      workTime.Date,
		})
	}

	return &result, nil
}

func (c *WorkTimeController) DeleteWorkTime(ctx context.Context,
	req *gen.DeleteWorkTimeRequestV1) (*emptypb.Empty, error) {

	params := entity.DeleteWorkTimeDTO{
		MasterID: req.GetMasterId(),
		Date:     req.GetDate(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &emptypb.Empty{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	err = c.deleteService.DeleteWorkTimeByMasterIDAndDate(ctx, params)
	if err != nil {
		return &emptypb.Empty{}, errors.HandleGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (c *WorkTimeController) UpdateWorkTime(ctx context.Context,
	req *gen.UpdateWorkTimeRequestV1) (*gen.WorkTimeEntityV1, error) {

	params := entity.UpdateWorkTimeDTO{
		MasterID:  req.GetMasterId(),
		Date:      req.GetDate(),
		StartTime: req.GetStartTime(),
		EndTime:   req.GetEndTime(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.WorkTimeEntityV1{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.updateService.UpdateWorkTime(ctx, params)
	if err != nil {
		return &gen.WorkTimeEntityV1{}, errors.HandleGrpcError(err)
	}
	return &gen.WorkTimeEntityV1{
		Id:        result.ID,
		MasterId:  result.MasterID,
		Date:      result.Date,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
	}, nil
}

func NewWorkTimeController(
	deleteService service.DeleteInteractor,
	getService service.GetInteractor,
	updateService service.UpdateInteractor,
	val *govalidator.Validator,
) *WorkTimeController {
	return &WorkTimeController{
		deleteService: deleteService,
		getService:    getService,
		updateService: updateService,
		val:           val,
	}
}
