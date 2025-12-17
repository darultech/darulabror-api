package service

import (
	"darulabror/internal/dto"
	"darulabror/internal/models"
	"darulabror/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService interface {
	// Superadmin only
	CreateAdmin(requesterRole models.Role, adminDTO dto.AdminDTO) error
	GetAllAdmins(page, limit int) ([]dto.AdminDTO, int64, error)
	UpdateAdmin(requesterRole models.Role, adminDTO dto.AdminDTO) error
	DeleteAdmin(requesterRole models.Role, id uint) error

	// shared (admin/superadmin)
	GetAdminByID(id uint) (dto.AdminDTO, error)

	// Public (login)
	AuthenticateAdmin(email, password string) (string, dto.AdminDTO, error)
}

type adminService struct {
	repo      repository.AdminRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAdminService(repo repository.AdminRepository, jwtSecret string) AdminService {
	return &adminService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtTTL:    24 * time.Hour,
	}
}

func (s *adminService) AuthenticateAdmin(email, password string) (string, dto.AdminDTO, error) {
	if len(s.jwtSecret) == 0 {
		return "", dto.AdminDTO{}, errors.New("JWT secret is not configured")
	}

	admin, err := s.repo.GetAdminByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", dto.AdminDTO{}, ErrInvalidCredentials
		}
		logrus.WithError(err).WithField("email", email).Error("failed get admin by email")
		return "", dto.AdminDTO{}, err
	}

	if !admin.IsActive {
		return "", dto.AdminDTO{}, ErrAdminInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", dto.AdminDTO{}, ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"admin_id": admin.ID,
		"role":     admin.Role,
		"iat":      now.Unix(),
		"exp":      now.Add(s.jwtTTL).Unix(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.jwtSecret)
	if err != nil {
		logrus.WithError(err).Error("failed sign jwt")
		return "", dto.AdminDTO{}, err
	}

	out := dto.AdminModelToDTO(admin)
	out.Password = "" // jangan expose hash
	return signed, out, nil
}

func (s *adminService) CreateAdmin(requesterRole models.Role, adminDTO dto.AdminDTO) error {
	if requesterRole != models.Superadmin {
		logrus.WithField("requester_role", requesterRole).Warn("forbidden create admin")
		return errors.New("forbidden")
	}

	// prevent duplicate email (simple check)
	if existing, err := s.repo.GetAdminByEmail(adminDTO.Email); err == nil && existing.ID != 0 {
		logrus.WithField("email", adminDTO.Email).Warn("admin email already exists")
		return ErrInvalidAdmin
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("failed checking admin email")
		return err
	}

	admin, err := dto.AdminDTOToModel(adminDTO)
	if err != nil {
		logrus.WithError(err).Error("failed to convert AdminDTO to model")
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("failed to hash admin password")
		return err
	}
	admin.Password = string(hash)

	if err := s.repo.CreateAdmin(admin); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"email": admin.Email,
			"role":  admin.Role,
		}).Error("failed to create admin")
		return ErrCreateAdmin
	}

	logrus.WithFields(logrus.Fields{
		"email": admin.Email,
		"role":  admin.Role,
	}).Info("admin created")
	return nil
}

func (s *adminService) GetAllAdmins(page, limit int) ([]dto.AdminDTO, int64, error) {
	admins, total, err := s.repo.GetAllAdmins(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed to get all admins")
		return nil, 0, err
	}

	out := make([]dto.AdminDTO, 0, len(admins))
	for _, a := range admins {
		d := dto.AdminModelToDTO(a)
		d.Password = "" // jangan expose hash/password ke response
		out = append(out, d)
	}

	return out, total, nil
}

func (s *adminService) GetAdminByID(id uint) (dto.AdminDTO, error) {
	admin, err := s.repo.GetAdminByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AdminDTO{}, ErrNotFoundAdmin
		}
		logrus.WithError(err).WithField("id", id).Error("failed to get admin by id")
		return dto.AdminDTO{}, err
	}

	d := dto.AdminModelToDTO(admin)
	d.Password = ""
	return d, nil
}

func (s *adminService) UpdateAdmin(requesterRole models.Role, adminDTO dto.AdminDTO) error {
	if requesterRole != models.Superadmin {
		logrus.WithField("requester_role", requesterRole).Warn("forbidden update admin")
		return errors.New("forbidden")
	}
	if adminDTO.ID == 0 {
		return ErrInvalidAdmin
	}

	admin, err := s.repo.GetAdminByID(adminDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFoundAdmin
		}
		return err
	}

	// update mutable fields
	admin.Username = adminDTO.Username
	admin.Email = adminDTO.Email
	admin.Role = adminDTO.Role
	if adminDTO.IsActive != nil {
		admin.IsActive = *adminDTO.IsActive
	}

	if adminDTO.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(adminDTO.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.WithError(err).Error("failed to hash admin password")
			return err
		}
		admin.Password = string(hash)
	}

	if err := s.repo.UpdateAdmin(admin); err != nil {
		logrus.WithError(err).WithField("id", admin.ID).Error("failed to update admin")
		return err
	}

	logrus.WithField("id", admin.ID).Info("admin updated")
	return nil
}

func (s *adminService) DeleteAdmin(requesterRole models.Role, id uint) error {
	if requesterRole != models.Superadmin {
		logrus.WithField("requester_role", requesterRole).Warn("forbidden delete admin")
		return errors.New("forbidden")
	}

	if err := s.repo.DeleteAdmin(id); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to delete admin")
		return err
	}

	logrus.WithField("id", id).Info("admin deleted")
	return nil
}
