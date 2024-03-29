package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"github.com/sonyamoonglade/delivery-service/pkg/cli"
	"github.com/sonyamoonglade/delivery-service/pkg/formatter"
	"github.com/sonyamoonglade/delivery-service/pkg/postgres"
	"github.com/sonyamoonglade/delivery-service/pkg/telegram"
	"github.com/sonyamoonglade/notification-service/pkg/logging"
)

func main() {
	log.Println("booting an application")

	logsPath, debug, strictMode := parseFlags()

	logger, err := logging.WithConfig(&logging.Config{
		Strict:   strictMode,
		LogsPath: logsPath,
		Debug:    debug,
		Encoding: logging.JSON,
	})

	if err != nil {
		log.Println(err.Error())
	}

	//Load .env.local for local development
	if err = godotenv.Load(".env"); err != nil {
		logger.Warnf("Ignore this message if app is ran by docker. %s", err.Error())
	}

	appCfg, err := config.GetAppConfig()
	if err != nil {
		logger.Fatalf("Could not read from config. %s", err.Error())
	}

	db, err := postgres.Connect(appCfg.Db)
	if err != nil {
		logger.Fatalf("Could not connect to database. %s", err.Error())
	}
	logger.Info("Database has connected")

	appBot, err := bot.NewBot(appCfg.Bot, logger)
	if err != nil {
		logger.Fatalf("Could not initialize newBot. %s", err.Error())
	}
	logger.Info("Bot has initialized")

	cliClient := cli.NewCli(logger)

	if err := cliClient.Ping(); err != nil {
		logger.Fatalf(err.Error())
	}

	extractFmt := formatter.NewFormatter(logger)

	//Initialize storage
	runnerStorage := runnStorage.NewRunnerStorage(db)
	deliveryStorage := dlvStorage.NewDeliveryStorage(db)

	//Initialize service
	checkService := check.NewCheckService(appCfg.App.CheckPath)
	deliveryService := dlvService.NewDeliveryService(logger, deliveryStorage, cliClient, checkService)
	runnerService := runnService.NewRunnerService(logger, runnerStorage)

	//Initialize transport
	telegramHandler := telegram.NewTelegramTransport(logger, appBot, runnerService, deliveryService, extractFmt)
	deliveryHandler := dlvHttp.NewDeliveryHandler(logger, deliveryService, extractFmt, appBot)
	runnerHandler := runnHttp.NewRunnerHandler(logger, runnerService)

	//Initialize router
	router := httprouter.New()

	deliveryHandler.RegisterRoutes(router)
	runnerHandler.RegisterRoutes(router)
	logger.Info("API Routes has initialized")

	go telegramHandler.ListenForUpdates()
	logger.Info("Bot is listening to updates")
	//Bot cant run more than 1 instance at a time
	//d := telegramHandler
	//_ = d
	server := tgdelivery.NewServerWithConfig(appCfg.App, router)

	//Start listening to requests
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server could not start. %s", err.Error())
		}
	}()
	logger.Infof("API server is listening on port %s", appCfg.App.Port)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)

	//Graceful shutdown
	<-exit
	logger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("could not shutdown httpserver. %s", err.Error())
	}
	logger.Info("ok")
}

func parseFlags() (string, bool, bool) {

	logsPath := flag.String("logs-path", "", "defines path to logging file")
	debug := flag.Bool("debug", true, "defines debug mode")
	strictMode := flag.Bool("strict", false, "defines strictness of the logs")

	flag.Parse()

	return *logsPath, *debug, *strictMode
}
