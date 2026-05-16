package jwt

type ConfigParser struct {
	Issuer    string `env:"ISSUER" env-required`
	Format    string `env:"FORMAT" env-required`
	Algorithm string `env:"ALGORITHM" env-required`
	Kid       string `env:"KID" env-required`
}
