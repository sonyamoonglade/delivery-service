package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/sonyamoonglade/delivery-service/pkg/bot"
	"github.com/sonyamoonglade/delivery-service/pkg/postgres"
	"github.com/spf13/viper"
)

type App struct {
	Port string
	Os   string
}

type AppConfig struct {
	Bot *bot.Config
	Db  *postgres.Config
	App *App
}

const (
	BotToken    = "BOT_TOKEN"
	BotUrl      = "BOT_URL"
	AdminUrl    = "ADMIN_URL"
	GroupChatID = "GROUP_CHAT_ID"
	DbPassword  = "DB_PASSWORD"
)

var TempOffset int64

func GetAppConfig() (AppConfig, error) {

	v, err := readConfig()
	if err != nil {
		return AppConfig{}, err
	}

	botToken, ok := os.LookupEnv(BotToken)
	if ok != true {
		return AppConfig{}, errors.New("missing botToken")
	}
	botURL, ok := os.LookupEnv(BotUrl)
	if ok != true {
		return AppConfig{}, errors.New("missing botURL")
	}

	adminURL, ok := os.LookupEnv(AdminUrl)
	if ok != true {
		return AppConfig{}, errors.New("missing adminURL")
	}

	groupChatID, ok := os.LookupEnv(GroupChatID)
	if ok != true {
		return AppConfig{}, errors.New("missing groupChatID")
	}

	chatIdLikeInt, err := strconv.ParseInt(groupChatID, 10, 64)
	if err != nil {
		return AppConfig{}, err
	}

	botCfg := &bot.Config{
		BotToken:    botToken,
		Timeout:     60, //idle timeout
		Debug:       false,
		URL:         botURL,
		AdminLink:   adminURL,
		GroupChatID: chatIdLikeInt,
	}

	dbUser := v.GetString("db.user")
	if dbUser == "" {
		return AppConfig{}, errors.New("missing db.user")
	}

	dbPwd, ok := os.LookupEnv(DbPassword)
	if ok != true {
		return AppConfig{}, errors.New("missing dbPassword")
	}

	dbHost := v.GetString("db.host")
	if dbHost == "" {
		return AppConfig{}, errors.New("missing db.host")
	}

	dbPort := v.GetString("db.port")
	if dbPort == "" {
		return AppConfig{}, errors.New("missing db.port")
	}

	dbPortLikeInt, err := strconv.ParseInt(dbPort, 10, 64)
	if err != nil {
		return AppConfig{}, err
	}

	dbName := v.GetString("db.database")
	if dbName == "" {
		return AppConfig{}, errors.New("missing db.database")
	}

	dbCfg := &postgres.Config{
		User:     dbUser,
		Password: dbPwd,
		Host:     dbHost,
		Port:     dbPortLikeInt,
		Database: dbName,
	}

	appPort := v.GetString("app.port")
	if appPort == "" {
		return AppConfig{}, errors.New("missing app.port")
	}

	baseOffset := v.GetString("app.baseOffset")
	if baseOffset == "" {
		return AppConfig{}, errors.New("missing app.baseOffset")
	}

	offsetLikeInt, err := strconv.ParseInt(baseOffset, 10, 64)
	if err != nil {
		return AppConfig{}, err
	}
	TempOffset = offsetLikeInt

	opSys, ok := os.LookupEnv("GOOS")
	if ok != true {
		return AppConfig{}, errors.New("missing GOOS")
	}

	appCfg := &App{
		Port: appPort,
		Os:   opSys,
	}

	return AppConfig{
		Bot: botCfg,
		Db:  dbCfg,
		App: appCfg,
	}, nil

}

func readConfig() (*viper.Viper, error) {

	env, ok := os.LookupEnv("ENV")
	if ok != true {
		return nil, errors.New("missing ENV")
	}

	name := "config"

	if env == "production" {
		name = "prod.config"
	}

	fmt.Printf("reading %s\n", name)

	viper.AddConfigPath(".")
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return viper.GetViper(), nil
}
