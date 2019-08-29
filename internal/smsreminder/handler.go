package smsreminder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type SmsReminderHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetOlderThan(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
	UpdateByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type smsReminderHandler struct {
	smsReminderService SmsReminderService
}

func NewSmsReminderHandler(smsReminderService SmsReminderService) SmsReminderHandler {
	return &smsReminderHandler{
		smsReminderService,
	}
}

func (h *smsReminderHandler) GetOlderThan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	time := vars["time"]

	fmt.Printf("handler.go getOlderThan() time=%v\n", time)

	smsReminders, err := h.smsReminderService.FindOlderThanSmsReminders(time)
	if err != nil {
		logrus.WithField("error", err).Error("Unable to find all smsReminders")
		http.Error(w, "Unable to find all smsReminders", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(smsReminders)
	if err != nil {
		logrus.WithField("error", err).Error("Error unmarshalling response")
		http.Error(w, "Unable to get smsReminder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		logrus.WithField("error", err).Error("Error writing response")
	}
}

func (h *smsReminderHandler) Get(w http.ResponseWriter, r *http.Request) {
	smsReminders, err := h.smsReminderService.FindAllSmsReminders()
	if err != nil {
		logrus.WithField("error", err).Error("Unable to find all smsReminders")
		http.Error(w, "Unable to find all smsReminders", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(smsReminders)
	if err != nil {
		logrus.WithField("error", err).Error("Error unmarshalling response")
		http.Error(w, "Unable to get smsReminder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		logrus.WithField("error", err).Error("Error writing response")
	}
}

func (h *smsReminderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	smsReminder, err := h.smsReminderService.FindSmsReminderByID(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "id": id}).Error("Unable to find smsReminder")
		http.Error(w, "Unable to find smsReminder", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(smsReminder)
	if err != nil {
		logrus.WithField("error", err).Error("Error unmarshalling response")
		http.Error(w, "Unable to fetch smsReminders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		logrus.WithField("error", err).Error("Error writing response")
	}
}

func (h *smsReminderHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {

}

func (h *smsReminderHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler.go DeleteByID() called")
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.smsReminderService.DeleteSmsReminderByID(id)
	if err != nil {
		logrus.WithField("error", err).Error("Error calling smsReminderService.DeleteSmsReminderByID")
		http.Error(w, "Unable to delete smsReminder", http.StatusInternalServerError)
		return
	}
}

func (h *smsReminderHandler) Create(w http.ResponseWriter, r *http.Request) {

	var smsReminder SmsReminder
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&smsReminder); err != nil {
		logrus.Error("Unable to decode smsReminder", err)
		http.Error(w, "Bad format for smsReminder", http.StatusBadRequest)
		return
	}

	if err := h.smsReminderService.CreateSmsReminder(&smsReminder); err != nil {
		logrus.WithField("error", err).Error("Unable to create smsReminder")
		http.Error(w, "Unable to create smsReminder", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(smsReminder)
	if err != nil {
		logrus.WithField("error", err).Error("Error unmarshalling response")
		http.Error(w, "Unable to create smsReminder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(response); err != nil {
		logrus.WithField("error", err).Error("Error writing response")
	}
}
