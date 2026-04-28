package jwt

type jwtTokenParser struct {
	jwtCache   JwtCache
	authClient // gRPC
}

func NewJwtTokenParser(jwtCache JwtCache) *jwtTokenParser {
	return &jwtTokenParser{
		jwtCache: jwtCache,
	}
}
