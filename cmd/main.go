package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"github.com/sonyamoonglade/delivery-service/config"
	apihandler "github.com/sonyamoonglade/delivery-service/internal/handler/http_api"
	tghandler "github.com/sonyamoonglade/delivery-service/internal/handler/tg_api"
	"github.com/sonyamoonglade/delivery-service/internal/service"
	"github.com/sonyamoonglade/delivery-service/internal/storage"
	"go.uber.org/zap"
	"os"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		logger.Error(err.Error())
	}

	if err = godotenv.Load(".env"); err != nil {
		logger.Error("Could not load environment variables")
	}

	appCfg, err := config.ReadConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Could not read from config. %s", err.Error()))
	}

	db, err := storage.Connect(&storage.PostgresConfig{
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
	logger.Info("Telegram update handler initialized")

	go tgHandler.ListenForUpdates(bot, updCfg)
	logger.Info("Bot is listening to updates")

	tgservice := service.NewTelegramService(logger, bot)

	deliveryStorage := storage.NewDeliveryStorage(logger, db)
	deliveryService := service.NewDeliveryService(logger, deliveryStorage)
	deliveryHandler := apihandler.NewDeliveryHandler(logger, deliveryService, tgservice)
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
