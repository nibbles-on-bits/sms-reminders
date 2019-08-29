package smsreminder

import (
	"time"
)

type SmsReminder struct {
	ID            string    `json:"id" db:"id"`
	FromNumber    string    `json:"fromNumber" db:"from_number"`
	ToNumber      string    `json:"toNumber" db:"to_number"`
	Message       string    `json:"message" db:"message"`
	ScheduledTime time.Time `json:"scheduledTime" db:"scheduled_time"`
	CreatedTime   time.Time `json:"createdTime" db:"created_time"`
	UpdatedTime   time.Time `json:"updatedTime" db:"updated_time"`
	DeletedTime   time.Time `json:"deletedTime" db:"deleted_time"`
}
