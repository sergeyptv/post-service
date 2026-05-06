package domain

type User struct {
	Uuid         string
	Username     string
	PasswordHash string
	Email        string
}
