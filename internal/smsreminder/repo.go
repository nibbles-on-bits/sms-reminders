package smsreminder

type SmsReminderRepository interface {
	Create(smsReminder *SmsReminder) error
	FindAll() ([]*SmsReminder, error)
	FindByID(id string) (*SmsReminder, error)
	DeleteByID(id string) error
}
