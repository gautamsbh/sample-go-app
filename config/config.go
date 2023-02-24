package config

type Config struct {
	Host string
	Port int
}

var AppConfig Config

// initialize config
//
// To read from environment variable, os.Getenv("ENV_NAME")
func init() {
	AppConfig = Config{
		Host: "",
		Port: 8000,
	}
}
