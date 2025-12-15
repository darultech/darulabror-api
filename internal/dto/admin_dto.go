package dto

import "darulabror/internal/models"

type AdminDTO struct {
	ID        uint        `json:"id" validate:"omitempty"`
	Username  string      `json:"username" validate:"required,min=3,max=50"`
	Email     string      `json:"email" validate:"required,email"`
	Password  string      `json:"password" validate:"required,min=6,max=50"`
	Role      models.Role `json:"role" validate:"required,oneof=admin superadmin;default=admin"`
	IsActive  bool        `json:"is_active" validate:"required"`
	CreatedAt int64       `json:"created_at,omitempty"`
	UpdatedAt int64       `json:"updated_at,omitempty"`
}

func AdminDTOToModel(dto AdminDTO) (models.Admin, error) {
	return models.Admin{
		ID:       dto.ID,
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
		IsActive: dto.IsActive,
	}, nil
}

func AdminModelToDTO(admin models.Admin) AdminDTO {
	return AdminDTO{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Password:  admin.Password,
		Role:      admin.Role,
		IsActive:  admin.IsActive,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}
