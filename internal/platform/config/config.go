package config

type App struct {
	Name    string `env:"NAME" env-prefix:"APP_" env-required`
	Version string `env:"VERSION" env-prefix:"APP_" env-required`
	Env     string `env:"ENV" env-prefix:"APP_" env-required`
}
