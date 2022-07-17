package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/config"
	dlvService "github.com/sonyamoonglade/delivery-service/internal/delivery/service"
	dlvStorage "github.com/sonyamoonglade/delivery-service/internal/delivery/storage"
	dlvHttp "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/http"
	runnService "github.com/sonyamoonglade/delivery-service/internal/runner/service"
	runnStorage "github.com/sonyamoonglade/delivery-service/internal/runner/storage"
	runnHttp "github.com/sonyamoonglade/delivery-service/internal/runner/transport/http"
	tgService "github.com/sonyamoonglade/delivery-service/internal/telegram/service"
	tgTransport "github.com/sonyamoonglade/delivery-service/internal/telegram/transport"
	"github.com/sonyamoonglade/delivery-service/pkg/postgres"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {

	logger, err := zap.NewProduction()

	if err != nil {
		log.Println(err.Error())
	}

	if err = godotenv.Load(".env"); err != nil {
		logger.Error("Could not load environment variables")
	}

	appCfg, err := config.ReadConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Could not read from config. %s", err.Error()))
	}

	db, err := postgres.Connect(&postgres.DbConfig{
		User:     appCfg.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     appCfg.GetString("db.host"),
		Port:     appCfg.GetInt64("db.port"),
		Database: appCfg.GetString("db.database"),
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Could not connect to database. %s", err.Error()))
	}
	logger.Info("Database has connected")

	botCfg := &tgdelivery.BotConfig{
		Token:   os.Getenv(tgdelivery.BOT_TOKEN),
		Timeout: 60,
		Debug:   false,
	}
	bot, updCfg, err := tgdelivery.BotWithConfig(botCfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not initialize bot. %s", err.Error()))
	}
	logger.Info("Bot has initialized")

	//Initialize storage
	runnerStorage := runnStorage.NewRunnerStorage(db)
	deliveryStorage := dlvStorage.NewDeliveryStorage(db)

	//Initialize service
	deliveryService := dlvService.NewDeliveryService(logger, deliveryStorage)
	telegramService := tgService.NewTelegramService(logger, bot)
	runnerService := runnService.NewRunnerService(logger, runnerStorage)

	//Initialize transport
	telegramHandler := tgTransport.NewTgHandler(logger, bot, runnerService)
	deliveryHandler := dlvHttp.NewDeliveryHandler(logger, deliveryService, telegramService)
	runnerHandler := runnHttp.NewRunnerHandler(logger, runnerService)

	//Initialize router
	router := httprouter.New()

	deliveryHandler.RegisterRoutes(router)
	runnerHandler.RegisterRoutes(router)
	logger.Info("API Routes has initialized")

	go telegramHandler.ListenForUpdates(bot, updCfg)
	logger.Info("Bot is listening to updates")

	server := tgdelivery.NewServerWithConfig(appCfg, router)
	logger.Info(fmt.Sprintf("API server is listening on port %s", appCfg.GetString("app.port")))
	if err = server.ListenAndServe(); err != nil {
		logger.Error(fmt.Sprintf("Server could not start. %s", err.Error()))
	}

}
