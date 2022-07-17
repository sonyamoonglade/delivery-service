package runner

import "github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"

type Storage interface {
	IsRunner(dto dto.IsRunnerDto) (bool, error)
	Register(dto dto.RegisterRunnerDto) (int64, error)
	Ban(id int64) (int64, error)
}
