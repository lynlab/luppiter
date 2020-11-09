package models

type UserIdentity struct {
	ModelMixin
	UUID     string
	Username string
	Email    string
	Accounts []UserAccount
}
