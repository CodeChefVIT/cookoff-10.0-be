package utils

import (
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type cfg struct {
	Port             string `env:"PORT" envDefault:"8080"`
	JwtSecret        string `env:"JWT_SECRET,notEmpty"`
	PostgresHost     string `env:"POSTGRES_HOST,notEmpty"`
	PostgresPort     string `env:"POSTGRES_PORT,notEmpty"`
	PostgresUser     string `env:"POSTGRES_USER,notEmpty"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,notEmpty"`
	PostgresDB       string `env:"POSTGRES_DB,notEmpty"`
	RedisHost        string `env:"REDIS_HOST,notEmpty"`
	RedisPort        string `env:"REDIS_PORT,notEmpty"`
	RedisPassword    string `env:"REDIS_PASSWORD,notEmpty"`
	SecretKey        string `env:"SECRET_KEY,notEmpty"`
	CallbackURL      string `env:"CALLBACK_URL,notEmpty"`
	Judge0URI        string `env:"JUDGE0_URI,notEmpty"`
	CookieSecure     bool   `env:"SECURE" envDefault:"false"`
}

var Config cfg

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		logger.Errorf("No .env file found")
	}

	if err := env.Parse(&Config); err != nil {
		logger.Errorf("%+v", err)
	}

	logger.Infof("Configuration successfully loaded")
}
