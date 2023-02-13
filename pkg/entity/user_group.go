package entity

type UserGroup struct {
	Base

	Name                 string   `json:"name" gorm:"unique;not null"`
	Organization         string   `json:"organization" gorm:"not null"`
	InvitationCodeDigest string   `json:"-" gorm:"not null"`
	IsFullAccess         bool     `json:"is_full_access" gorm:"not null"`
	Bastion              *Bastion `json:"bastion,omitempty"`
}
