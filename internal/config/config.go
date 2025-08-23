package config

import ( 
	"time"
	"os"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 		string 		`yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string 		`yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer	`yaml:"http_server"`
}

type HTTPServer struct {
	Adderss 	string 			`yaml:"address" env-default"localhost:8080"`
	Timeout 	time.Duration 	`yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration 	`yaml:"idle_timeout" env-default:"60s"`
}
//"F:/Rest_api/config/local.yaml"
func MustLoad() *Config {
	configPath := os.Getenv("CONF_PATH")
	if configPath == "" {
		log.Fatal("CONF_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}