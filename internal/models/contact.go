package models

type ContactStatus string

const (
	ContactStatusNew        ContactStatus = "new"
	ContactStatusInProgress ContactStatus = "in_progress"
	ContactStatusDone       ContactStatus = "done"
)

type Contact struct {
	ID        uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string        `gorm:"not null" json:"email"`
	Subject   string        `gorm:"not null" json:"subject"`
	Message   string        `gorm:"type:text;not null" json:"message"`
	Status    ContactStatus `gorm:"type:text;not null;default:'new';check:status IN ('new','in_progress','done')" json:"status"`
	CreatedAt int64         `gorm:"autoCreateTime" json:"created_at"`
}
