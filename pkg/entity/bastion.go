package entity

import "github.com/google/uuid"

type Bastion struct {
	Base

	UserGroupID uuid.UUID `json:"user_group_id" gorm:"not null"`
	User        string    `json:"bastion_user" gorm:"not null"`
	Password    string    `json:"bastion_password" gorm:"not null"`
	Host        string    `json:"bastion_host" gorm:"not null"`
	Port        int       `json:"bastion_port" gorm:"not null"`
}
