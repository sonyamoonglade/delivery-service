package runner

import (
	"github.com/sonyamoonglade/delivery-service/internal/entity"
	"github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
)

type Service interface {
	IsRunner(usrPhoneNumber string) (int64, error)
	IsKnownByTelegramId(tgUsrID int64) (bool, error)
	GetByTelegramId(tgUsrID int64) (*entity.Runner, error)
	Register(dto dto.RegisterRunnerDto) error
	Ban(runnerID int64) error
	BeginWork(dto dto.RunnerBeginWorkDto) error
}
