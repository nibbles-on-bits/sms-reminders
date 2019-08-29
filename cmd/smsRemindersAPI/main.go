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
	DefaultHTTPPort      = "3003"
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

	TwilioSID := env.EnvString("SRM_TWILIO_SID", "")
	TwilioToken := env.EnvString("SRM_TWILIO_TOKEN", "")
	HTTPPort := env.EnvString("SRM_HTTP_PORT", DefaultHTTPPort)

	fmt.Printf("TwilioSID = %s\n", TwilioSID)
	fmt.Printf("TwilioToken = %s\n", TwilioToken)

	if TwilioSID == "" {
		Logger.Panic("environment variable : TWILIO_SID not set")
	}

	if TwilioToken == "" {
		Logger.Panic("environment variable : TWILIO_TOKEN not set")
	}
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
		statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS sms_reminders (id TEXT PRIMARY KEY, from_number TEXT, to_number TEXT, message TEXT, scheduled_time INTEGER, created_time INTEGER, updated_time INTEGER, deleted_time INTEGER, processing BOOL)")
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
	router.HandleFunc("/smsreminders/olderthan/{time}", SmsReminderHandler.GetOlderThan).Methods("GET")
	router.HandleFunc("/smsreminders/{id}", SmsReminderHandler.DeleteByID).Methods("DELETE")
	router.HandleFunc("/smsreminders", SmsReminderHandler.Create).Methods("POST")
	//router.HandleFunc("/events")	// TODO : an event receiver endpoint for Event Sourcing

	errs := make(chan error, 2)

	go func() {
		logrus.Info(fmt.Sprintf("Listening server mode on port : %s", HTTPPort))
		p := ":" + HTTPPort
		errs <- http.ListenAndServe(p, router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("Err Chan %s", <-c)
	}()

	/*go func() {
		for true {
			fmt.Println("Timer")
			time.Sleep(5 * time.Second)

			// given the current time stamp
		}
	}()*/

	logrus.Error("sms-reminders-microservice terminated", <-errs)

}

/*func sendSMS(sr SmsReminder) {
	sr.FromNumber
	sr.ToNumber

}
*/
// type SmsReminder struct {
// 	ID            string    `json:"id" db:"id"`
// 	FromNumber    string    `json:"fromNumber" db:"from_number"`
// 	ToNumber      string    `json:"toNumber" db:"to_number"`
// 	Message       string    `json"message" db:"message"`
// 	ScheduledTime time.Time `json:"scheduledTime" db:"scheduled_time"`
// 	CreatedTime   time.Time `json:"createdTime" db:"created_time"`
// 	UpdatedTime   time.Time `json:"updatedTime" db:"updated_time"`
// 	DeletedTime   time.Time `json:"deletedTime" db:"deleted_time"`
// }

//TODO : Create a smsRemindersService, smsRemindersHandler
