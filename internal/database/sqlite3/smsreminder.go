package sqlite3

import (
	"database/sql"
	"fmt"
	"log"
	"sms-reminders-microservice/internal/smsreminder"
	"strconv"
	"time"
)

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

	tmScheduled := strconv.FormatInt(sr.ScheduledTime.Unix(), 10)
	tmCreated := strconv.FormatInt(sr.CreatedTime.Unix(), 10)
	tmUpdated := strconv.FormatInt(sr.UpdatedTime.Unix(), 10)
	tmDeleted := strconv.FormatInt(sr.DeletedTime.Unix(), 10)

	res, err := statement.Exec(sr.ID, sr.FromNumber, sr.ToNumber, sr.Message, tmScheduled, tmCreated, tmUpdated, tmDeleted)

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
	rows, err := r.db.Query("SELECT id, from_number, to_number, message, scheduled_time, created_Time, updated_time, deleted_time FROM sms_reminders")
	defer rows.Close()

	for rows.Next() {
		sr := new(smsreminder.SmsReminder)
		// tmScheduled := ""
		// tmCreated := ""
		// tmUpdated := ""
		// tmDeleted := ""

		var tmScheduled int64
		var tmCreated int64
		var tmUpdated int64
		var tmDeleted int64

		if err = rows.Scan(&sr.ID, &sr.FromNumber, &sr.ToNumber, &sr.Message, &tmScheduled, &tmCreated, &tmUpdated, &tmDeleted); err != nil {
			log.Print(err)
			return nil, err
		}

		fmt.Println("Debugging FindAll()")
		fmt.Printf("tmScheduled=%d\n", tmScheduled)

		t := time.Unix(tmScheduled, 0)
		fmt.Println(t, err)
		sr.ScheduledTime = t
		t = time.Unix(tmCreated, 0)

		fmt.Println(t, err)
		sr.CreatedTime = t
		t = time.Unix(tmUpdated, 0)

		fmt.Println(t, err)
		sr.UpdatedTime = t
		t = time.Unix(tmDeleted, 0)

		fmt.Println(t, err)
		sr.DeletedTime = t

		smsReminders = append(smsReminders, sr)

	}
	return smsReminders, nil
}

// DeleteByID attempts to delete a smsReminder in a sqlite3 repository
func (r *smsReminderRepository) DeleteByID(id string) error {
	_, err := r.db.Exec("DELETE FROM sms_reminders where id=$1", id)

	if err != nil {
		panic(err)
	}

	return nil
}
