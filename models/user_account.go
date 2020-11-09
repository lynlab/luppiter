package models

type UserAccount struct {
	ModelMixin
	Provider   string
	ProviderID string
	IdentityID int64
	Identity   UserIdentity
}
