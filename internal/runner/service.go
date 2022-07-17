package runner

import "github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"

type Service interface {
	IsRunner(dto dto.IsRunnerDto) (bool, error)
	Register(dto dto.RegisterRunnerDto) error
	Ban(id int64) error
}
