package ports

type TokenSigner interface {
	NewToken(userUuid, username, userEmail string) (string, error)
}
