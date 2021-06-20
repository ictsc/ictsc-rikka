package entity

type UserGroup struct {
	Base

	Name                 string `json:"name"`
	Organization         string `json:"organization"`
	InvitationCodeDigest string `json:"-"`
	IsFullAccess         bool   `json:"is_full_access"`
}
