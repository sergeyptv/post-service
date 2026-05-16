package config

type App struct {
	Name    string `env:"NAME" env-required`
	Version string `env:"VERSION" env-required`
	Env     string `env:"ENV" env-required`
}
