package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"os/exec"
	"strings"
	"sync"
)

var TimeoutError = errors.New("operation takes too long")
var CliError = errors.New("internal cli error")

type Cli interface {
	WriteCheck(dto dto.CheckDtoForCli) error
	Ping() error
}

type cli struct {
	logger *zap.SugaredLogger
	mut    sync.Mutex
}

func NewCli(logger *zap.SugaredLogger) Cli {
	return &cli{
		logger: logger,
		mut:    sync.Mutex{},
	}
}

func (c *cli) WriteCheck(dto dto.CheckDtoForCli) error {
	//mutex here

	c.mut.Lock()
	defer c.mut.Unlock()
	byt, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	strForCli := string(byt)

	var stdErr buffer.Buffer

	command := fmt.Sprintf("bin/cli.exe")

	// pass -dto flag with string dto
	cmd := exec.Command(command, "-dto", fmt.Sprintf(`%s`, strForCli))

	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		//If error occurs -> return
		errText := strings.ToLower(stdErr.String())
		if strings.Contains(errText, "api key has expired") {
			return check.ApiKeyHasExpired
		}
		c.logger.Errorf("CLI call error. stderr: %s", errText)
		return CliError
	}

	//Command has run successfully
	c.logger.Info("CLI call has been successful")

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
		//If error occurs -> parse stdErr to normal err

		return err
	}

	//Command has run successfully
	c.logger.Info("CLI ping has been successful")

	return nil
}
