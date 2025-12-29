package service

import (
	"darulabror/internal/models"
	"darulabror/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContactService interface {
	// Public
	CreateContact(email, subject, message string) error

	// Admin
	GetAllContacts(page, limit int, status string) ([]models.Contact, int64, error)
	GetContactByID(id uint) (*models.Contact, error)
	UpdateContact(id uint, email, subject, message string) error
	UpdateContactStatus(id uint, status models.ContactStatus) error
	DeleteContact(id uint) error
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) CreateContact(email, subject, message string) error {
	if err := s.repo.CreateContact(email, subject, message); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"email":   email,
			"subject": subject,
		}).Error("failed create contact")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"email":   email,
		"subject": subject,
	}).Info("contact created")
	return nil
}

func (s *contactService) GetAllContacts(page, limit int, status string) ([]models.Contact, int64, error) {
	contacts, total, err := s.repo.GetAllContacts(page, limit, status)
	if err != nil {
		logrus.WithError(err).Error("failed get all contacts")
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *contactService) GetContactByID(id uint) (*models.Contact, error) {
	contact, err := s.repo.GetContactByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		logrus.WithError(err).WithField("id", id).Error("failed get contact by id")
		return nil, err
	}
	return contact, nil
}

func (s *contactService) UpdateContact(id uint, email, subject, message string) error {
	if err := s.repo.UpdateContact(id, email, subject, message); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed update contact")
		return err
	}
	logrus.WithField("id", id).Info("contact updated")
	return nil
}

func (s *contactService) DeleteContact(id uint) error {
	if err := s.repo.DeleteContact(id); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed delete contact")
		return err
	}
	logrus.WithField("id", id).Info("contact deleted")
	return nil
}

func (s *contactService) UpdateContactStatus(id uint, status models.ContactStatus) error {
	// Validate status value
	if status != models.ContactStatusNew && status != models.ContactStatusInProgress && status != models.ContactStatusDone {
		return errors.New("invalid status value")
	}
	
	if err := s.repo.UpdateContactStatus(id, status); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed update contact status")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"id":     id,
		"status": status,
	}).Info("contact status updated")
	return nil
}
