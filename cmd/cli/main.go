package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"github.com/sonyamoonglade/delivery-service/pkg/logging"
	"go.uber.org/zap"
	"log"
	"os"
)

const keysPath = "check/keys.txt"
const templatePath = "check/check_template.docx"
const keyHasRestored = "key has restored"

func main() {

	log.Println("booting check-formatter app")

	logger, err := logging.WithCfg(&logging.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		DevMode:  false,
		Encoding: logging.JSON,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	flags := RegisterFlags()
	//var inp dto.CheckDto
	flag.Parse()
	logger.Debug("parsed flags")
	flagDto := *flags.dto

	var cliDto dto.CheckDtoForCli

	if err := json.Unmarshal([]byte(flagDto), &cliDto); err != nil {
		log.Fatal(err.Error())
	}
	logger.Debug("got cli dto from command line")

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

	key, err := check.GetFirstKey()
	if err != nil {
		if errors.Is(err, check.NoApiKeysLeft) {
			logger.Error(err.Error())
			return
		}
		logger.Fatal(err.Error())
		return
	}
	logger.Debugf("obtained key %s", key)

	if err := check.SetLicense(key); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("set a license")

	template, err := check.OpenTemplate(templatePath)
	if err != nil {
		//Api key is no longer valid
		if errors.Is(err, check.ApiKeyHasExpired) {
			//Restore key and exit
			if err := check.RestoreKey(); err != nil {
				//Zero api keys has left. Exit
				if errors.Is(err, check.NoApiKeysLeft) {
					logger.Fatal(err.Error())
					return
				}
				//Internal error
				logger.Fatal(err.Error())
				return
			}
			//Successfully restored a key
			logger.Info(keyHasRestored)
			return
		}
		//Internal error
		logger.Fatal(err.Error())
		return
	}

	// '{"data":{"order":{"order_id":1,"total_cart_price":1080,"pay":"withCard","cart":[{"name":"Моцарелла","price":499,"quantity":2}],"is_delivered":false},"user":{"username":"Иван Семенов","phone_number":"+79128507000"}}}'
	check.Format(template, inp)
	logger.Info("formatted a template")

	if err := template.SaveToFile("./check/check.docx"); err != nil {
		logger.Fatal(err.Error())
		return
	}

	logger.Info("successfully saved check.docx")
	return
}

func RegisterFlags() Flags {
	return Flags{
		dto: flag.String("dto", "", "dto"),
	}

}

type Flags struct {
	dto *string
}
