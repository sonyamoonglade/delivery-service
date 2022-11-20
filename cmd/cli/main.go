package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"github.com/sonyamoonglade/notification-service/pkg/logging"
)

const keysPath = "check/keys.txt"
const templatePath = "check/check_template.docx"

func main() {

	var cliDto dto.CheckDtoForCli

	path := "check" // Transformed into check/file.ext
	checkService := check.NewCheckService(path)

	log.Println("booting check-formatter app")

	logger, err := logging.WithConfig(&logging.Config{
		Strict:   false, //hardcode
		Debug:    true,  //hardcode
		LogsPath: "",    //hardcode
		Encoding: logging.JSON,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	//Parse and register flags
	flags := RegisterFlags()
	flag.Parse()

	//Allow clients to ping cli
	ping := *flags.ping
	flagDto := *flags.dto
	if ping == true {
		return
	}

	if err := json.Unmarshal([]byte(flagDto), &cliDto); err != nil {
		log.Fatal(err.Error(), "json error")
	}

	inp := cliDto.Data

	//Check if keys file exists
	if _, err := os.Stat(keysPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Fatal(err.Error())
			return
		}
		logger.Fatalf(err.Error())
		return
	}
	//Check if template file exists
	if _, err := os.Stat(templatePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Fatal(err.Error())
			return
		}
		logger.Fatal(err.Error())
		return
	}
	logger.Info("all files are ok")

	//Get first key from keys
	key, err := checkService.GetFirstKey()
	if err != nil {
		if errors.Is(err, check.NoApiKeysLeft) {
			logger.Fatal(err.Error())
			return
		}
		logger.Fatalf("caused check.GetFirstKey. %s", err.Error())
		return
	}
	logger.Debug("obtained key", key)

	if err := checkService.SetLicense(key); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("set a license")

	//Open docx template
	template, err := checkService.OpenTemplate(templatePath)
	if err != nil {
		//Api key is no longer valid
		if errors.Is(err, check.ApiKeyHasExpired) {
			logger.Fatal(err.Error())
			return
		}
		if errors.Is(err, check.NoApiKeysLeft) {
			logger.Fatalf("caused check.OpenTemplate. %s", err.Error())
			return
		}
		//Internal error
		logger.Fatal(err.Error())
		return
	}

	//Format template with input from command line
	checkService.Format(template, inp)
	logger.Info("formatted a template")

	//Save formatted template
	if err := template.SaveToFile("check/check.docx"); err != nil {
		logger.Fatal(err.Error())
		return
	}
	defer template.Close()

	logger.Info("successfully saved check.docx")
	return
}

func RegisterFlags() Flags {
	return Flags{
		ping: flag.Bool("ping", false, "ping"),
		dto:  flag.String("dto", "", "dto"),
	}

}

type Flags struct {
	dto  *string
	ping *bool
}
