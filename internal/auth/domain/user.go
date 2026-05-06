package domain

type User struct {
	Uuid     string
	Username string
	PassHash string
	Email    string
}
