package domain

type TokenSigner interface {
	NewToken(userUuid, userEmail string) (string, error)
}
