package repository

import (
	"darulabror/internal/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	//Manage Admins by Superadmin
	CreateAdmin(admin models.Admin) error
	GetAllAdmins() ([]models.Admin, error)
	GetAdminByUserID(userID string) (models.Admin, error)
	UpdateAdmin(admin models.Admin) error
	DeleteAdmin(id uint) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateAdmin(admin models.Admin) error {
	return r.db.Create(&admin).Error
}

func (r *adminRepository) GetAllAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	err := r.db.Find(&admins).Error
	return admins, err
}

func (r *adminRepository) GetAdminByUserID(userID string) (models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("user_id = ?", userID).First(&admin).Error
	return admin, err
}

func (r *adminRepository) UpdateAdmin(admin models.Admin) error {
	return r.db.Save(&admin).Error
}

func (r *adminRepository) DeleteAdmin(id uint) error {
	return r.db.Delete(&models.Admin{}, id).Error
}
