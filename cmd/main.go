package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kachmazoff/doit-api"
	_ "github.com/kachmazoff/doit-api/docs"
	"github.com/kachmazoff/doit-api/internal/controller"
	"github.com/kachmazoff/doit-api/internal/mailing"
	"github.com/kachmazoff/doit-api/internal/repository"
	"github.com/kachmazoff/doit-api/internal/service"
)

// @title Course Platform API
// @version 1.0
// @description API Server for Course Platform

// @host localhost:8000
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey StudentsAuth
// @in header
// @name Authorization

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	smtpConfig := mailing.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	var sender mailing.Sender
	if smtpConfig.Username == "" || smtpConfig.Password == "" {
		log.Println("could not setup smtp mailing: username and/or password not provided; using log implementation")
		sender = mailing.NewLogSender(
			smtpConfig.From,
			&logWrapper{Logger: log.New(os.Stdout, "mailing ", log.LstdFlags)},
		)
	} else {
		sender = mailing.NewSMTPSender(smtpConfig)
	}

	db, err := repository.NewMysqlDB(repository.Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_Port"),
		Username: os.Getenv("DB_Username"),
		Password: os.Getenv("DB_Password"),
		DBName:   os.Getenv("DB_Name"),
		SSLMode:  os.Getenv("DB_SSLMode"),
	})

	if err != nil {
		log.Fatalf("Failed to init db: %s", err.Error())
	}

	port := env("PORT", "8080")

	repos := repository.NewMysqlRepos(db)
	services := service.NewServices(repos, sender)
	controllers := controller.NewController(services)

	// user := model.User{
	// 	Username: "root",
	// 	Email:    "alek.kachmazov@yandex.ru",
	// 	Password: "root",
	// }

	// userId, _ := services.Users.Create(user)
	// services.Users.ConfirmAccount(userId)
	// return

	srv := new(doit.Server)
	go func() {
		if err := srv.Run(port, controllers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("Doit API Started")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Doit API Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error occured on db connection close: %s", err.Error())
	}
}

func env(key, fallbackValue string) string {
	s, ok := os.LookupEnv(key)
	if !ok || s == "" {
		return fallbackValue
	}
	return s
}

type logWrapper struct {
	Logger *log.Logger
}

func (l *logWrapper) Log(args ...interface{}) {
	l.Logger.Println(args...)
}

func (l *logWrapper) Logf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
}
