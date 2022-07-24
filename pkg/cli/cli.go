package cli

import (
	"encoding/json"
	"fmt"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"os/exec"
)

type Cli interface {
	WriteCheck(dto dto.CheckDtoForCli) error
	Ping() error
}

type cli struct {
	logger *zap.SugaredLogger
}

func NewCli(logger *zap.SugaredLogger) Cli {
	return &cli{
		logger: logger,
	}
}

func (c *cli) WriteCheck(dto dto.CheckDtoForCli) error {

	byt, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	strForCli := fmt.Sprintf("'%s'", string(byt))

	var stdOut buffer.Buffer
	var stdErr buffer.Buffer

	command := fmt.Sprintf("bin/cli.exe -dto %s", strForCli)
	c.logger.Debugf("command: %s", command)

	cmd := exec.Command(command)

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if err := cmd.Run(); err != nil {
		//If error occurs -> return
		c.logger.Errorf("stdErr: %s", stdErr.String())
		return err
	}

	//Command has run successfully
	c.logger.Info("CLI call has been successful")
	c.logger.Infof("CLI stdout: %s", stdOut.String())

	return nil
}

func (c *cli) Ping() error {

	c.logger.Info("pinging cli")

	var stdOut buffer.Buffer
	var stdErr buffer.Buffer

	command := "bin/cli.exe"
	flags := "-ping"
	c.logger.Debugf("command: %s. flags: %s", command, flags)

	cmd := exec.Command(command, flags)

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if err := cmd.Run(); err != nil {
		//If error occurs -> return
		c.logger.Errorf("stdErr: %s", stdErr.String())
		return err
	}

	//Command has run successfully
	c.logger.Info("CLI ping has been successful")

	return nil
}
