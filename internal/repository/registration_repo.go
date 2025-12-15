package repository

import (
	"darulabror/internal/models"
	"errors"

	"gorm.io/gorm"
)

type RegistrationRepo interface {
	// Public Registration Management
	Create(reg models.Registration) error
	// Admin Registration Management
	GetAll() ([]models.Registration, error)
	GetByID(id uint) (models.Registration, error)
	GetByEmail(email string) (models.Registration, error)
	GetByNISN(nisn string) (models.Registration, error)

	Update(reg models.Registration) error
	Delete(id uint) error
	// Existence Checks
	ExistsByEmail(email string) (bool, error)
	ExistsByNISN(nisn string) (bool, error)
}

type registrationRepo struct {
	db *gorm.DB
}

func NewRegistrationRepo(db *gorm.DB) RegistrationRepo {
	return &registrationRepo{db: db}
}

func (r *registrationRepo) Create(reg models.Registration) error {
	return r.db.Create(&reg).Error
}

func (r *registrationRepo) GetAll() ([]models.Registration, error) {
	var regs []models.Registration
	err := r.db.Order("id DESC").Find(&regs).Error
	return regs, err
}

func (r *registrationRepo) GetByID(id uint) (models.Registration, error) {
	var reg models.Registration
	err := r.db.First(&reg, id).Error
	return reg, err
}

func (r *registrationRepo) GetByEmail(email string) (models.Registration, error) {
	var reg models.Registration
	err := r.db.Where("email = ?", email).First(&reg).Error
	return reg, err
}

func (r *registrationRepo) GetByNISN(nisn string) (models.Registration, error) {
	var reg models.Registration
	err := r.db.Where("nisn = ?", nisn).First(&reg).Error
	return reg, err
}

func (r *registrationRepo) Update(reg models.Registration) error {
	// Pastikan ID ada
	if reg.ID == 0 {
		return errors.New("registration id is required")
	}
	return r.db.Save(&reg).Error
}

func (r *registrationRepo) Delete(id uint) error {
	return r.db.Delete(&models.Registration{}, id).Error
}

func (r *registrationRepo) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Registration{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *registrationRepo) ExistsByNISN(nisn string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Registration{}).Where("nisn = ?", nisn).Count(&count).Error
	return count > 0, err
}
