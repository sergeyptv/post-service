package jwt

type ConfigParser struct {
	Issuer    string `env:"ISSUER" env-prefix:"TOKEN_" env-required`
	Format    string `env:"FORMAT" env-prefix:"TOKEN_" env-required`
	Algorithm string `env:"ALGORITHM" env-prefix:"TOKEN_" env-required`
}
