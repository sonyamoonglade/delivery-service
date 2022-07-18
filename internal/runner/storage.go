package runner

import "github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"

type Storage interface {
	IsRunner(usrPhoneNumber string) (int64, error)
	IsKnownByTelegramId(usrID int64) (bool, error)
	GetByTelegramId(tgUsrID int64) (int64, error)
	Register(dto dto.RegisterRunnerDto) (int64, error)
	Ban(runnerID int64) (int64, error)
	BeginWork(dto dto.RunnerBeginWorkDto) error
}
