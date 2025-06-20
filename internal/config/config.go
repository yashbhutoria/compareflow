package config

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AllowedOrigins []string
}