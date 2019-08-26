package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"sms-reminders-microservice/internal/database/sqlite3"
	"sms-reminders-microservice/internal/env"
	"sms-reminders-microservice/internal/smsreminder"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

const (
	DefaultRedisURL      = "localhost:6379"
	DefaultRedisPassword = "example"

	DefaultPostgresURL      = "postgres://example:postgres@localhost/vehicle?sslmode=disable"
	DefaultPostgresUser     = "postgres"
	DefaultPostgresHost     = "db"
	DefaultPostgresPort     = "5432"
	DefaultPostgresPassword = "example"
	DefaultPostgresDBName   = "smsreminders_microservice_db"

	DefaultSqlite3File = "./smsreminders.db"
)

func main() {
	errChan := make(chan error)

	Logger := zap.NewExample()
	defer Logger.Sync()
	Logger.Info("Welcome to sms-reminders-microservice")

	var SmsReminderRepo smsreminder.SmsReminderRepository

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	var dbURL string
	dbType := "sqlite3" // for now, we can add others later

	switch dbType {
	case "sqlite3":
		dbURL = env.EnvString("DATABASE_URL", DefaultSqlite3File)
		db, err := sql.Open("sqlite3", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS sms_reminders (id TEXT PRIMARY KEY, from_number TEXT, to_number TEXT, message TEXT, scheduled_time TEXT, created_time TEXT, updated_time TEXT, deleted_time TEXT)")
		statement.Exec()
		defer db.Close()
		SmsReminderRepo = sqlite3.NewSqlite3SmsReminderRepository(db)
	default:
		panic("Unknown database")
	}

	SmsReminderService := smsreminder.NewSmsReminderService(SmsReminderRepo)
	SmsReminderHandler := smsreminder.NewSmsReminderHandler(SmsReminderService)

	router := mux.NewRouter()
	router.HandleFunc("/smsreminders", SmsReminderHandler.Get).Methods("GET")
	router.HandleFunc("/smsreminders/{id}", SmsReminderHandler.GetByID).Methods("GET")
	router.HandleFunc("/smsreminders/{id}", SmsReminderHandler.DeleteByID).Methods("DELETE")
	router.HandleFunc("/smsreminders", SmsReminderHandler.Create).Methods("POST")
	//router.HandleFunc("/events")	// TODO : an event receiver endpoint for Event Sourcing

	errs := make(chan error, 2)

	go func() {
		logrus.Info("Listening server mode on port :3003")
		//errs <- http.ListenAndServe(":7777", nil)
		errs <- http.ListenAndServe(":3003", router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("Err Chan %s", <-c)
	}()

	logrus.Error("sms-reminders-microservice terminated", <-errs)

}

//TODO : Create a smsRemindersService, smsRemindersHandler
