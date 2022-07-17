package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/config"
	service2 "github.com/sonyamoonglade/delivery-service/internal/delivery/service"
	storage2 "github.com/sonyamoonglade/delivery-service/internal/delivery/storage"
	apihandler "github.com/sonyamoonglade/delivery-service/internal/delivery/transport/http"
	"github.com/sonyamoonglade/delivery-service/internal/telegram/service"
	tghandler "github.com/sonyamoonglade/delivery-service/internal/telegram/transport"
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

	tgHandler := tghandler.NewTgHandler(logger, bot)
	tgService := service.NewTelegramService(logger, bot)
	logger.Info("Telegram composite initialized")

	go tgHandler.ListenForUpdates(bot, updCfg)
	logger.Info("Bot is listening to updates")

	deliveryStorage := storage2.NewDeliveryStorage(logger, db)
	deliveryService := service2.NewDeliveryService(logger, deliveryStorage)
	deliveryHandler := apihandler.NewDeliveryHandler(logger, deliveryService, tgService)
	logger.Info("Delivery composite initialized")

	router := httprouter.New()
	deliveryHandler.RegisterRoutes(router)
	logger.Info("API Routes has initialized")

	server := tgdelivery.NewServerWithConfig(appCfg, router)
	logger.Info(fmt.Sprintf("API server is listening on port %s", appCfg.GetString("app.port")))
	if err = server.ListenAndServe(); err != nil {
		logger.Error(fmt.Sprintf("Server could not start. %s", err.Error()))
	}

}
