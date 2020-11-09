package models

type Application struct {
	ModelMixin
	UUID      string
	Name      string
	OwnerID   int
	Owner     UserIdentity
	SecretKey string
}
