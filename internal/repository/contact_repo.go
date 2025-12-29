package repository

import (
	"darulabror/internal/models"
	"darulabror/internal/utils"

	"gorm.io/gorm"
)

type ContactRepository interface {
	// Public methods for contact Admin
	CreateContact(email, subject, message string) error
	// Admin methods for contact
	GetAllContacts(page, limit int, status string) ([]models.Contact, int64, error)
	GetContactByID(id uint) (*models.Contact, error)
	UpdateContact(id uint, email, subject, message string) error
	UpdateContactStatus(id uint, status models.ContactStatus) error
	DeleteContact(id uint) error
}

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) CreateContact(email, subject, message string) error {
	return r.db.Create(&models.Contact{
		Email:   email,
		Subject: subject,
		Message: message,
		Status:  models.ContactStatusNew, // default status
	}).Error
}

func (r *contactRepository) GetAllContacts(page, limit int, status string) ([]models.Contact, int64, error) {
	var (
		contacts []models.Contact
		total    int64
	)

	_, limit, offset := utils.NormalizePageLimit(page, limit)

	query := r.db.Model(&models.Contact{})
	
	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&contacts).Error
	return contacts, total, err
}

func (r *contactRepository) GetContactByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, id).Error
	return &contact, err
}

func (r *contactRepository) UpdateContact(id uint, email, subject, message string) error {
	var contact models.Contact
	if err := r.db.First(&contact, id).Error; err != nil {
		return err
	}

	contact.Email = email
	contact.Subject = subject
	contact.Message = message
	return r.db.Save(&contact).Error
}

func (r *contactRepository) DeleteContact(id uint) error {
	return r.db.Delete(&models.Contact{}, id).Error
}

func (r *contactRepository) UpdateContactStatus(id uint, status models.ContactStatus) error {
	result := r.db.Model(&models.Contact{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
