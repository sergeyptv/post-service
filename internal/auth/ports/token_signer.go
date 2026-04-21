package ports

type TokenSigner interface {
	NewToken(userUuid, userEmail string) (string, error)
}
