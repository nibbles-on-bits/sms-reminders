package smsreminder

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// SmsReminderService is the interface for smsReminders
type SmsReminderService interface {
	CreateSmsReminder(smsReminder *SmsReminder) error
	FindSmsReminderByID(id string) (*SmsReminder, error)
	FindAllSmsReminders() ([]*SmsReminder, error)
	DeleteSmsReminderByID(id string) error
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

func (s *smsReminderService) DeleteSmsReminderByID(id string) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "id": id}).Error("Error deleting smsReminder")
		return err
	}
	return err
}
