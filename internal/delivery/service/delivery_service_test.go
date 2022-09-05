package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	mock_delivery "github.com/sonyamoonglade/delivery-service/internal/delivery/mocks"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	mock_check "github.com/sonyamoonglade/delivery-service/pkg/check/mocks"
	"github.com/sonyamoonglade/delivery-service/pkg/cli"
	mock_cli "github.com/sonyamoonglade/delivery-service/pkg/cli/mocks"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/test/global"
	"github.com/stretchr/testify/require"
)

func initServiceDeps(ctrl *gomock.Controller) (*mock_delivery.MockStorage, *mock_cli.MockCli, *mock_check.MockService) {
	cliMock := mock_cli.NewMockCli(ctrl)
	checkServiceMock := mock_check.NewMockService(ctrl)
	deliveryStorage := mock_delivery.NewMockStorage(ctrl)
	return deliveryStorage, cliMock, checkServiceMock
}

func TestStatusOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.StatusOfDeliveryDto{
		OrderIDs: []int64{1, 2, 3},
	}

	mockStorageResp := []bool{true, false, true}

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Status(inp.OrderIDs).Return(mockStorageResp, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	statuses, err := deliveryService.Status(inp)
	require.NoError(t, err)
	require.Equal(t, len(statuses), len(inp.OrderIDs))

	for i, status := range statuses {
		expID := inp.OrderIDs[i]
		require.Equal(t, expID, status.OrderID)
		require.Equal(t, mockStorageResp[i], status.Status)
	}

}
func TestStatusErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.StatusOfDeliveryDto{
		OrderIDs: []int64{1, 2, 3},
	}

	mockErr := errors.New("some crazy database error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Status(inp.OrderIDs).Return(nil, mockErr).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	statuses, err := deliveryService.Status(inp)
	require.Error(t, err)
	require.Nil(t, statuses)

	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
}

func TestCreateOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CreateDeliveryDatabaseDto{
		OrderID: 2,
		Pay:     "onPickup",
	}
	//DeliveryID returned
	mockResp := int64(2)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Create(inp).Return(mockResp, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	deliveryID, err := deliveryService.Create(inp)
	require.NoError(t, err)
	require.Equal(t, mockResp, deliveryID)

}
func TestCreateAlreadyExists(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CreateDeliveryDatabaseDto{
		OrderID: 2,
		Pay:     "onPickup",
	}
	//DeliveryID returned already exists - 0
	mockResp := int64(0)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Create(inp).Return(mockResp, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	deliveryID, err := deliveryService.Create(inp)
	require.Equal(t, mockResp, deliveryID)
	require.Error(t, err)
	require.Equal(t, httpErrors.ConflictError(httpErrors.DeliveryAlreadyExists).Error(), err.Error())
}
func TestCreateErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CreateDeliveryDatabaseDto{
		OrderID: 2,
		Pay:     "onPickup",
	}
	//DeliveryID returned already exists - 0
	mockErr := errors.New("some database error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Create(inp).Return(int64(0), mockErr).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	deliveryID, err := deliveryService.Create(inp)
	require.Equal(t, int64(0), deliveryID)
	require.Error(t, err)
	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
}

func TestCompleteOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := int64(2)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Complete(inp).Return(true, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	ok, err := deliveryService.Complete(inp)
	require.Equal(t, true, ok)
	require.NoError(t, err)
}
func TestCompleteErr(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := int64(2)
	mockErr := errors.New("crazy db error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Complete(inp).Return(false, mockErr).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	ok, err := deliveryService.Complete(inp)
	require.Equal(t, false, ok)
	require.Error(t, err)
	require.Equal(t, err.Error(), err.Error())
}
func TestCompleteErrDeliveryCantBeCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := int64(2)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Complete(inp).Return(false, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	ok, err := deliveryService.Complete(inp)
	require.Equal(t, false, ok)
	require.Error(t, err)
	require.Equal(t, tgErrors.DeliveryCouldNotBeCompleted(inp).Error(), err.Error())
}

func TestReserveOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.ReserveDeliveryDto{
		RunnerID:   2,
		DeliveryID: 3,
	}

	mockReturn := time.Now()

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Reserve(inp).Return(mockReturn, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	reservedAt, err := deliveryService.Reserve(inp)
	require.NoError(t, err)
	require.Equal(t, time.Now(), reservedAt)
	require.NotZero(t, reservedAt)
}
func TestReserveErr(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.ReserveDeliveryDto{
		RunnerID:   2,
		DeliveryID: 3,
	}

	mockErr := errors.New("crazy db err")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Reserve(inp).Return(time.Time{}, mockErr).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	reservedAt, err := deliveryService.Reserve(inp)
	require.Error(t, err)
	require.Equal(t, mockErr.Error(), err.Error())
	require.Zero(t, reservedAt)
	require.Equal(t, time.Time{}, reservedAt)

}
func TestReserveAlreadyReserved(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.ReserveDeliveryDto{
		RunnerID:   2,
		DeliveryID: 3,
	}

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Reserve(inp).Return(time.Time{}, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	reservedAt, err := deliveryService.Reserve(inp)
	require.Error(t, err)
	require.Zero(t, reservedAt)
	require.Equal(t, tgErrors.DeliveryHasAlreadyReserved(inp.DeliveryID).Error(), err.Error())
}

func TestDeleteOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := int64(4)

	//Dep
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Delete(inp).Return(true, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.Delete(inp)
	require.NoError(t, err)
	require.Nil(t, err)
}
func TestDeleteNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := int64(4)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Delete(inp).Return(false, nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.Delete(inp)
	require.Error(t, err)
	require.NotNil(t, err)
	require.Equal(t, httpErrors.NotFoundError(httpErrors.DeliveryDoesNotExist).Error(), err.Error())
}
func TestDeleteErr(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := int64(2)
	mockErr := errors.New("crazy error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Delete(inp).Return(false, mockErr).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.Delete(inp)
	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
	require.Error(t, err)
	require.NotNil(t, err)
}

func TestWriteCheckOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CheckDtoForCli{
		Data: dto.CheckDto{
			Order: tgdelivery.OrderForCheck{
				OrderID: 24,
				DeliveryDetails: tgdelivery.DeliveryDetails{
					Address:        "Test",
					FlatCall:       1,
					EntranceNumber: 2,
					Floor:          3,
					DeliveredAt:    time.Now(),
					Comment:        "test comment",
				},
				Amount: 25,
				Pay:    "onPickup",
				Cart: []tgdelivery.CartProduct{
					{
						ProductID: 1,
						Name:      "Mozzarella",
						Price:     200,
						Quantity:  1,
						Category:  "Пицца",
					},
				},
				IsDelivered: true,
			},
			User: tgdelivery.UserForCheck{
				Username:    "Bobby Martin",
				PhoneNumber: "+79128572849",
			},
		},
	}

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	cliMock.EXPECT().WriteCheck(inp).Return(nil).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.WriteCheck(inp)
	require.NoError(t, err)
	require.Nil(t, err)
}
func TestWriteCheckRestoreKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CheckDtoForCli{
		Data: dto.CheckDto{
			Order: tgdelivery.OrderForCheck{
				OrderID: 24,
				DeliveryDetails: tgdelivery.DeliveryDetails{
					Address:        "Test",
					FlatCall:       1,
					EntranceNumber: 2,
					Floor:          3,
					DeliveredAt:    time.Now(),
					Comment:        "test comment",
				},
				Amount: 25,
				Pay:    "onPickup",
				Cart: []tgdelivery.CartProduct{
					{
						ProductID: 1,
						Name:      "Mozzarella",
						Price:     200,
						Quantity:  1,
						Category:  "Пицца",
					},
				},
				IsDelivered: true,
			},
			User: tgdelivery.UserForCheck{
				Username:    "Bobby Martin",
				PhoneNumber: "+79128572849",
			},
		},
	}

	//Deps
	deliveryStorageMock, cliMock, checkServiceMock := initServiceDeps(ctrl)

	/*
		This testing setup is kinda wierd...
		cliMock.WriteCheck is called 2 times and both should return ApiKeyHasExpired.
		Why 2 times? For loop. (see impl.)

		Well on each ApiKeyHasExpired return - checkServiceMock should call RestoreKey and return nil,
		indicating that restoring has occurred nicely.

		deliveryService.WriteCheck would return no err after, meaning everything is fine. <- wierd part. Does test anything though...

		This test suite tests whether checkService.RestoreKey is called if cli returns ApiKeyHasExpired.
	*/

	cliMock.EXPECT().WriteCheck(inp).Return(check.ApiKeyHasExpired).Times(2)

	//Suite should call RestoreKey and receive success(no err(nil))
	checkServiceMock.EXPECT().RestoreKey().Return(nil).Times(2)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkServiceMock)

	err := deliveryService.WriteCheck(inp)
	require.NoError(t, err)
	require.Nil(t, err)
}
func TestWriteCheckNoApiKeysLeft(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CheckDtoForCli{
		Data: dto.CheckDto{
			Order: tgdelivery.OrderForCheck{
				OrderID: 24,
				DeliveryDetails: tgdelivery.DeliveryDetails{
					Address:        "Test",
					FlatCall:       1,
					EntranceNumber: 2,
					Floor:          3,
					DeliveredAt:    time.Now(),
					Comment:        "test comment",
				},
				Amount: 25,
				Pay:    "onPickup",
				Cart: []tgdelivery.CartProduct{
					{
						ProductID: 1,
						Name:      "Mozzarella",
						Price:     200,
						Quantity:  1,
						Category:  "Пицца",
					},
				},
				IsDelivered: true,
			},
			User: tgdelivery.UserForCheck{
				Username:    "Bobby Martin",
				PhoneNumber: "+79128572849",
			},
		},
	}

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	//Should be called once
	cliMock.EXPECT().WriteCheck(inp).Return(check.NoApiKeysLeft).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.WriteCheck(inp)
	require.Equal(t, check.NoApiKeysLeft.Error(), err.Error())
	require.Error(t, err)
}
func TestWriteCheckTimeoutErr(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CheckDtoForCli{
		Data: dto.CheckDto{
			Order: tgdelivery.OrderForCheck{
				OrderID: 24,
				DeliveryDetails: tgdelivery.DeliveryDetails{
					Address:        "Test",
					FlatCall:       1,
					EntranceNumber: 2,
					Floor:          3,
					DeliveredAt:    time.Now(),
					Comment:        "test comment",
				},
				Amount: 25,
				Pay:    "onPickup",
				Cart: []tgdelivery.CartProduct{
					{
						ProductID: 1,
						Name:      "Mozzarella",
						Price:     200,
						Quantity:  1,
						Category:  "Пицца",
					},
				},
				IsDelivered: true,
			},
			User: tgdelivery.UserForCheck{
				Username:    "Bobby Martin",
				PhoneNumber: "+79128572849",
			},
		},
	}

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	//Should be called once
	cliMock.EXPECT().WriteCheck(inp).Return(cli.TimeoutError).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.WriteCheck(inp)
	require.Equal(t, cli.TimeoutError.Error(), err.Error())
	require.Error(t, err)
}
func TestWriteCheckInternalErrRestoringKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := global.InitLogger()

	inp := dto.CheckDtoForCli{
		Data: dto.CheckDto{
			Order: tgdelivery.OrderForCheck{
				OrderID: 24,
				DeliveryDetails: tgdelivery.DeliveryDetails{
					Address:        "Test",
					FlatCall:       1,
					EntranceNumber: 2,
					Floor:          3,
					DeliveredAt:    time.Now(),
					Comment:        "test comment",
				},
				Amount: 25,
				Pay:    "onPickup",
				Cart: []tgdelivery.CartProduct{
					{
						ProductID: 1,
						Name:      "Mozzarella",
						Price:     200,
						Quantity:  1,
						Category:  "Пицца",
					},
				},
				IsDelivered: true,
			},
			User: tgdelivery.UserForCheck{
				Username:    "Bobby Martin",
				PhoneNumber: "+79128572849",
			},
		},
	}
	mockErr := errors.New("some restoring key error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	//Should be called once
	cliMock.EXPECT().WriteCheck(inp).Return(check.ApiKeyHasExpired).Times(1)

	//Should be called once an then method would return internal error
	checkMock.EXPECT().RestoreKey().Return(mockErr).Times(1)

	//Under test
	deliveryService := NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	err := deliveryService.WriteCheck(inp)
	require.Error(t, err)
	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
	require.NotNil(t, err)
}
