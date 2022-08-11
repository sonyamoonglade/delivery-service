package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sonyamoonglade/delivery-service/config"
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

var PathToExecutable string

type Cli interface {
	WriteCheck(dto dto.CheckDtoForCli) error
	Ping() error
}

type cli struct {
	logger *zap.SugaredLogger
	mut    sync.Mutex
}

func NewCli(logger *zap.SugaredLogger, config *config.App) Cli {
	os := config.Os

	if os == strings.ToLower("linux") {
		PathToExecutable = "./bin/cli"
	}
	if os == strings.ToLower("windows") {
		PathToExecutable = "bin/cli.exe"
	}

	return &cli{
		logger: logger,
		mut:    sync.Mutex{},
	}
}

func (c *cli) WriteCheck(dto dto.CheckDtoForCli) error {

	c.mut.Lock()
	defer c.mut.Unlock()

	byt, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	strForCli := string(byt)

	var stdErr buffer.Buffer
	//Optional
	var stdOut buffer.Buffer

	// pass -dto flag with string dto
	cmd := exec.Command(PathToExecutable, "-dto", fmt.Sprintf(`%s`, strForCli))

	cmd.Stderr = &stdErr
	//Optional
	cmd.Stdout = &stdOut

	if err := cmd.Run(); err != nil {
		//If error occurs -> return
		errText := strings.ToLower(stdErr.String())
		if strings.Contains(errText, "api key has expired") {
			return check.ApiKeyHasExpired
		}
		c.logger.Errorf("CLI call error. stderr: %s", errText)
		return CliError
	}
	//Optional stdout
	c.logger.Infof("stdout: %s", stdOut.String())

	//Command has run successfully
	c.logger.Info("CLI call has been successful")

	return nil

}

func (c *cli) Ping() error {

	c.logger.Info("pinging cli")

	var stdOut buffer.Buffer
	var stdErr buffer.Buffer

	flags := "-ping"
	c.logger.Debugf("command: %s. flags: %s", PathToExecutable, flags)

	cmd := exec.Command(PathToExecutable, flags)

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
