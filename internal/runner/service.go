package runner

import "github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"

type Service interface {
	IsRunner(usrPhoneNumber string) (int64, error)
	Register(dto dto.RegisterRunnerDto) error
	Ban(runnerID int64) error
	BeginWork(dto dto.RunnerBeginWorkDto) error
}
