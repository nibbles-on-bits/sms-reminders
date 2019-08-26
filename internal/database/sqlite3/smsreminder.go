package sqlite3

import (
	"database/sql"
	"fmt"
	"log"
	"sms-reminders-microservice/internal/smsreminder"
	"time"
)

const SmsReminderTable = "SmsReminder"

type smsReminderRepository struct {
	db *sql.DB
}

func NewSqlite3SmsReminderRepository(db *sql.DB) smsreminder.SmsReminderRepository {
	return &smsReminderRepository{
		db,
	}
}

func (r *smsReminderRepository) Create(sr *smsreminder.SmsReminder) error {

	statement, err := r.db.Prepare("INSERT INTO sms_reminders(id, from_number, to_number, message, scheduled_time, created_Time, updated_time, deleted_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")

	if err != nil {
		panic(err)
	}

	tmCreated := sr.CreatedTime.UTC().Format(time.RFC3339)
	tmUpdated := sr.UpdatedTime.UTC().Format(time.RFC3339)
	tmDeleted := sr.DeletedTime.UTC().Format(time.RFC3339)

	res, err := statement.Exec(sr.ID, sr.FromNumber, sr.ToNumber, sr.Message, sr.ScheduledTime, tmCreated, tmUpdated, tmDeleted)

	fmt.Printf("res=%#v\n", res)

	if err != nil {
		panic(err)
	}
	return nil
}

func (r *smsReminderRepository) FindByID(id string) (*smsreminder.SmsReminder, error) {
	smsReminder := new(smsreminder.SmsReminder)
	err := r.db.QueryRow("SELECT id, from_number, to_number, message, scheduled_time, created_Time, updated_time, deleted_time FROM sms_reminders where id=$1", id).Scan(&smsReminder.ID, &smsReminder.FromNumber, &smsReminder.ToNumber, &smsReminder.Message, &smsReminder.ScheduledTime, &smsReminder.CreatedTime, &smsReminder.UpdatedTime, &smsReminder.DeletedTime)
	if err != nil {
		panic(err)
	}
	return smsReminder, nil
}

func (r *smsReminderRepository) FindAll() (smsReminders []*smsreminder.SmsReminder, err error) {
	rows, err := r.db.Query("SELECT id, smsReminder_number, year, make, model, vin, created, updated, deleted FROM smsReminders")
	defer rows.Close()

	for rows.Next() {
		smsReminder := new(smsreminder.SmsReminder)
		tmCreated := ""
		tmUpdated := ""
		tmDeleted := ""

		if err = rows.Scan(&smsReminder.ID, &smsReminder.SmsReminderNumber, &smsReminder.Year, &smsReminder.Make, &smsReminder.Model, &smsReminder.VIN, &tmCreated, &tmUpdated, &tmDeleted); err != nil {
			log.Print(err)
			return nil, err
		}

		t, err := time.Parse(time.RFC3339Nano, tmCreated)
		fmt.Println(t, err)
		smsReminder.Created = t
		t, err = time.Parse(time.RFC3339Nano, tmUpdated)
		fmt.Println(t, err)
		smsReminder.Updated = t
		t, err = time.Parse(time.RFC3339Nano, tmDeleted)
		fmt.Println(t, err)
		smsReminder.Deleted = t

		smsReminders = append(smsReminders, smsReminder)

	}
	return smsReminders, nil
}

// DeleteByID attempts to delete a smsReminder in a sqlite3 repository
func (r *smsReminderRepository) DeleteByID(id string) error {
	_, err := r.db.Exec("DELETE FROM smsReminders where id=$1", id)

	if err != nil {
		panic(err)
	}

	return nil
}
