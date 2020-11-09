package models

import (
	"time"
)

type AccessToken struct {
	ModelMixin

	IdentityID    int64
	Identity      UserIdentity
	ApplicationID int64
	Application   Application

	AccessKey     string
	SecretKey     string
	ActivationKey string
	Activated     bool

	ExpireAt *time.Time
}

func (t *AccessToken) HasExpired() bool {
	return t.ExpireAt.Before(time.Now())
}
