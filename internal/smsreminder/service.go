package smsreminder

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// SmsReminderService is the interface for smsReminders
type SmsReminderService interface {
	// Mark Reminder as expired
	// Find out if we have any reminders due
	//FindDueReminders(t time.Time, smsReminder *Sms) {}

	CreateSmsReminder(smsReminder *SmsReminder) error
	FindSmsReminderByID(id string) (*SmsReminder, error)
	FindAllSmsReminders() ([]*SmsReminder, error)
	FindOlderThanSmsReminders(time string) ([]*SmsReminder, error) // time should be unix time format 1567035250 would be Wednesday, August 28, 2019 11:34:10 PM Zulu
	DeleteSmsReminderByID(id string) error
	FindDueSmsReminders() ([]*SmsReminder, error)
}

type smsReminderService struct {
	repo SmsReminderRepository
}

// NewSmsReminderService will bind a repository
func NewSmsReminderService(repo SmsReminderRepository) SmsReminderService {
	return &smsReminderService{
		repo,
	}
}

func (s *smsReminderService) CreateSmsReminder(smsReminder *SmsReminder) error {
	smsReminder.ID = uuid.New().String()
	smsReminder.CreatedTime = time.Now()
	smsReminder.UpdatedTime = time.Now()

	if err := s.repo.Create(smsReminder); err != nil {
		logrus.WithField("error", err).Error("Error creating smsReminder")
		return err
	}

	logrus.WithField("id", smsReminder.ID).Info("Created new smsReminder")
	return nil
}

// FindOlderThanSmsReminders returns a collection of SmsReminders that are older than a particular zulu unix timestamp
func (s *smsReminderService) FindOlderThanSmsReminders(ut string) ([]*SmsReminder, error) {
	smsReminder, err := s.repo.FindAll()

	// Here we need to filter the ones older than t

	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Error finding all smsReminder older than : " + ut)
		return nil, err
	}
	logrus.Info("Found all smsReminders")
	return smsReminder, nil
}

func (s *smsReminderService) FindSmsReminderByID(id string) (*SmsReminder, error) {
	vehicle, err := s.repo.FindByID(id)

	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "id": id}).Error("Error finding smsReminder")
		return nil, err
	}
	logrus.WithField("id", id).Info("Found smsReminder")
	return vehicle, nil
}

func (s *smsReminderService) FindAllSmsReminders() ([]*SmsReminder, error) {
	smsReminder, err := s.repo.FindAll()
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Error finding all smsReminder")
		return nil, err
	}
	logrus.Info("Found all smsReminders")
	return smsReminder, nil
}

func (s *smsReminderService) FindDueSmsReminders() ([]*SmsReminder, error) {
	smsReminder, err := s.repo.FindAll()
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Error finding all smsReminder")
		return nil, err
	}
	logrus.Info("Found all smsReminders")
	return smsReminder, nil
}

func (s *smsReminderService) DeleteSmsReminderByID(id string) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "id": id}).Error("Error deleting smsReminder")
		return err
	}
	return err
}
