package entity

type UserGroup struct {
	Base

	Name                 string
	Organization         string
	InvitationCodeDigest string `json:"-"`
	IsFullAccess         bool
}
