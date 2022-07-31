package main

import (
	"context"
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
	bot "github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/check"
	"github.com/sonyamoonglade/delivery-service/pkg/cli"
	"github.com/sonyamoonglade/delivery-service/pkg/logging"
	"github.com/sonyamoonglade/delivery-service/pkg/postgres"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	log.Println("booting an application")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)

	logger, err := logging.WithCfg(&logging.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		DevMode:  true,
		Encoding: logging.JSON,
	})

	if err != nil {
		log.Println(err.Error())
	}

	//Load .env
	if err = godotenv.Load(); err != nil {
		logger.Error("Could not load environment variables")
	}

	appCfg, err := config.ReadConfig()
	if err != nil {
		logger.Fatalf("Could not read from config. %s", err.Error())
	}

	db, err := postgres.Connect(&postgres.DbConfig{
		User:     appCfg.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     appCfg.GetString("db.host"),
		Port:     appCfg.GetInt64("db.port"),
		Database: appCfg.GetString("db.database"),
	})
	if err != nil {
		logger.Fatalf("Could not connect to database. %s", err.Error())
	}
	logger.Info("Database has connected")

	grpChatID, err := strconv.ParseInt(os.Getenv("GROUP_CHAT_ID"), 10, 64)
	if err != nil {
		logger.Errorf("could not get group chat id. %s", err.Error())
	}

	botCfg := &bot.Config{
		Token:        os.Getenv(tgdelivery.BOT_TOKEN),
		Timeout:      60,
		Debug:        false,
		TelegramLink: os.Getenv("BOT_URL"),
		AdminLink:    os.Getenv("ADMIN_URL"),
		GroupChatID:  grpChatID,
	}
	botInstance, updCfg, err := bot.WithConfig(botCfg)
	if err != nil {
		logger.Fatalf("Could not initialize bot. %s", err.Error())
	}
	logger.Info("Bot has initialized")

	cliClient := cli.NewCli(logger)

	if err := cliClient.Ping(); err != nil {
		logger.Error(err.Error())
	}

	//Initialize storage
	runnerStorage := runnStorage.NewRunnerStorage(db)
	deliveryStorage := dlvStorage.NewDeliveryStorage(db)

	//Initialize service
	checkService := check.NewCheckService()
	deliveryService := dlvService.NewDeliveryService(logger, deliveryStorage, cliClient, checkService)
	telegramService := tgService.NewTelegramService(logger, botInstance)
	runnerService := runnService.NewRunnerService(logger, runnerStorage)

	//Initialize transport
	telegramHandler := tgTransport.NewTgHandler(logger, botInstance, runnerService, deliveryService, telegramService)
	deliveryHandler := dlvHttp.NewDeliveryHandler(logger, deliveryService, telegramService)
	runnerHandler := runnHttp.NewRunnerHandler(logger, runnerService)

	//Initialize router
	router := httprouter.New()

	deliveryHandler.RegisterRoutes(router)
	runnerHandler.RegisterRoutes(router)
	logger.Info("API Routes has initialized")

	go telegramHandler.ListenForUpdates(botInstance, updCfg)
	logger.Info("Bot is listening to updates")

	server := tgdelivery.NewServerWithConfig(appCfg, router)

	//Start listening to requests
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server could not start. %s", err.Error())
		}
	}()
	logger.Infof("API server is listening on port %s", appCfg.GetString("app.port"))

	//Graceful shutdown
	<-exit
	logger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	go func() {
		//todo
		defer cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("could not shutdown httpserver. %s", err.Error())
	}
	logger.Info("ok")
}
