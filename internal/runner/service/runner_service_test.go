package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	mock_runner "github.com/sonyamoonglade/delivery-service/internal/runner/mocks"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/test/global"
	"github.com/stretchr/testify/require"
)

//todo: move ctrl to tests..
func initDeps(t *testing.T) *mock_runner.MockStorage {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	runnerStorageMock := mock_runner.NewMockStorage(ctrl)
	return runnerStorageMock
}

func TestRegisterOK(t *testing.T) {

	logger := global.InitLogger()

	inp := dto.RegisterRunnerDto{
		Username:    "Alex Alex",
		PhoneNumber: "+79128523812",
	}
	mockReturn := int64(5)

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().Register(inp).Return(mockReturn, nil).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	err := runnerService.Register(inp)
	require.NoError(t, err)
	require.Nil(t, err)

}
func TestRegisterAlreadyExists(t *testing.T) {

	logger := global.InitLogger()

	inp := dto.RegisterRunnerDto{
		Username:    "Alex Alex",
		PhoneNumber: "+79128523812",
	}
	mockReturn := int64(0)

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().Register(inp).Return(mockReturn, nil).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	err := runnerService.Register(inp)
	require.Error(t, err)
	require.Equal(t, httpErrors.ConflictError(httpErrors.RunnerAlreadyExists).Error(), err.Error())
	require.NotNil(t, err)
}
func TestRegisterErr(t *testing.T) {

	logger := global.InitLogger()

	inp := dto.RegisterRunnerDto{
		Username:    "Alex Alex",
		PhoneNumber: "+79128523812",
	}
	mockError := errors.New("crazy err")

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().Register(inp).Return(int64(0), mockError).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	err := runnerService.Register(inp)
	require.Error(t, err)
	require.NotNil(t, err)
	require.Equal(t, httpErrors.InternalError().Error(), err.Error())
}

func TestIsKnownByTelegramOK(t *testing.T) {

	logger := global.InitLogger()

	inp := int64(-512312354)

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().IsKnownByTelegramId(inp).Return(true, nil).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	ok, err := runnerService.IsKnownByTelegramId(inp)
	require.NoError(t, err)
	require.Equal(t, true, ok)
	require.Nil(t, err)
}
func TestIsKnownByTelegramBadScenario(t *testing.T) {

	logger := global.InitLogger()

	inp := int64(-512312354)

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().IsKnownByTelegramId(inp).Return(false, nil).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	ok, err := runnerService.IsKnownByTelegramId(inp)
	require.NoError(t, err)
	require.Equal(t, false, ok)
	require.Nil(t, err)
}
func TestIsKnownByTelegramErr(t *testing.T) {

	logger := global.InitLogger()

	inp := int64(-512312354)
	mockErr := errors.New("crazy err")

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().IsKnownByTelegramId(inp).Return(false, mockErr).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	ok, err := runnerService.IsKnownByTelegramId(inp)
	require.Error(t, err)
	require.False(t, ok)
	require.NotNil(t, err)
	require.Equal(t, mockErr.Error(), err.Error())
}

func TestIsRunnerOK(t *testing.T) {
	logger := global.InitLogger()

	inp := "+79128509000"
	mockID := int64(5)

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().IsRunner(inp).Return(mockID, nil).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	runnerID, err := runnerService.IsRunner(inp)
	require.Equal(t, mockID, runnerID)
	require.Nil(t, err)
	require.NoError(t, err)
}
func TestIsRunnerDoesNotExist(t *testing.T) {
	logger := global.InitLogger()

	inp := "+79128509000"
	mockID := int64(0)

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().IsRunner(inp).Return(mockID, nil).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	runnerID, err := runnerService.IsRunner(inp)
	require.Equal(t, mockID, runnerID)
	require.NotNil(t, err)
	require.Error(t, err)
	require.Equal(t, tgErrors.RunnerDoesNotExist(inp).Error(), err.Error())
}
func TestIsRunnerErr(t *testing.T) {
	logger := global.InitLogger()

	inp := "+79128509000"
	mockID := int64(0)
	mockErr := errors.New("crazy err")

	//Deps
	runnerStorage := initDeps(t)

	runnerStorage.EXPECT().IsRunner(inp).Return(mockID, mockErr).Times(1)

	runnerService := NewRunnerService(logger, runnerStorage)

	runnerID, err := runnerService.IsRunner(inp)
	require.Equal(t, mockID, runnerID)
	require.NotNil(t, err)
	require.Equal(t, mockErr.Error(), err.Error())
}
