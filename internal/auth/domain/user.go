package domain

type User struct {
	Uuid     string
	Username string
	PassHash string
	Email    string
	Phone    string
}

type CreateUser struct {
	Username string
	PassHash string
	Email    string
	Phone    string
}
