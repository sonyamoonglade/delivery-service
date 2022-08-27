package delivery_service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	mock_delivery "github.com/sonyamoonglade/delivery-service/internal/delivery/mocks"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/service"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	mock_check "github.com/sonyamoonglade/delivery-service/pkg/check/mocks"
	mock_cli "github.com/sonyamoonglade/delivery-service/pkg/cli/mocks"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func initServiceDeps(ctrl *gomock.Controller) (*mock_delivery.MockStorage, *mock_cli.MockCli, *mock_check.MockService) {
	cliMock := mock_cli.NewMockCli(ctrl)
	checkServiceMock := mock_check.NewMockService(ctrl)
	deliveryStorage := mock_delivery.NewMockStorage(ctrl)
	return deliveryStorage, cliMock, checkServiceMock
}

func initLogger() *zap.SugaredLogger {
	prod, _ := zap.NewProduction()
	logger := prod.Sugar()
	return logger
}

func TestStatusOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

	inp := dto.StatusOfDeliveryDto{
		OrderIDs: []int64{1, 2, 3},
	}

	mockStorageResp := []bool{true, false, true}

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Status(inp.OrderIDs).Return(mockStorageResp, nil).Times(1)

	//Under test
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

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

	logger := initLogger()

	inp := dto.StatusOfDeliveryDto{
		OrderIDs: []int64{1, 2, 3},
	}

	mockErr := errors.New("some crazy database error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Status(inp.OrderIDs).Return(nil, mockErr).Times(1)

	//Under test
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	statuses, err := deliveryService.Status(inp)
	require.Error(t, err)
	require.Nil(t, statuses)

	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
}

func TestCreateOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

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
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	deliveryID, err := deliveryService.Create(inp)
	require.NoError(t, err)
	require.Equal(t, mockResp, deliveryID)

}
func TestCreateAlreadyExists(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

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
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	deliveryID, err := deliveryService.Create(inp)
	require.Equal(t, mockResp, deliveryID)
	require.Error(t, err)
	require.Equal(t, httpErrors.ConflictError(httpErrors.DeliveryAlreadyExists).Error(), err.Error())
}
func TestCreateErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

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
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	deliveryID, err := deliveryService.Create(inp)
	require.Equal(t, int64(0), deliveryID)
	require.Error(t, err)
	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
}

func TestCompleteOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

	inp := int64(2)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Complete(inp).Return(true, nil).Times(1)

	//Under test
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	ok, err := deliveryService.Complete(inp)
	require.Equal(t, true, ok)
	require.NoError(t, err)
}
func TestCompleteErr(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

	inp := int64(2)
	mockErr := errors.New("crazy db error")

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Complete(inp).Return(false, mockErr).Times(1)

	//Under test
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	ok, err := deliveryService.Complete(inp)
	require.Equal(t, false, ok)
	require.Error(t, err)
	require.Equal(t, err.Error(), err.Error())
}
func TestCompleteErrDeliveryCantBeCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := initLogger()

	inp := int64(2)

	//Deps
	deliveryStorageMock, cliMock, checkMock := initServiceDeps(ctrl)

	deliveryStorageMock.EXPECT().Complete(inp).Return(false, nil).Times(1)

	//Under test
	deliveryService := service.NewDeliveryService(logger, deliveryStorageMock, cliMock, checkMock)

	ok, err := deliveryService.Complete(inp)
	require.Equal(t, false, ok)
	require.Error(t, err)
	require.Equal(t, tgErrors.DeliveryCouldNotBeCompleted(inp).Error(), err.Error())
}
