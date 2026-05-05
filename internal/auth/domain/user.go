package domain

type User struct {
	Uuid     string
	Username string
	PassHash string
	Email    string
}

type InputUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
