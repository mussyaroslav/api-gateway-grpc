package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"`
	Server       `yaml:"server"`
	LogFile      `yaml:"logFile"`
	Hub          ClientGRPC `yaml:"hub"`
	Director     ClientGRPC `yaml:"director"`
	AuthRolesHub ClientGRPC `yaml:"authRolesHub"`
	Maestro      ClientGRPC `yaml:"maestro"`
	Certs        Certs      `yaml:"certs"`
	SSO          SSO        `yaml:"sso"`
}

type Server struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env-default:"8080"`
}

type LogFile struct {
	Use  bool   `yaml:"use" env-default:"false"`
	Name string `yaml:"name" env-default:"inc-director.log"`
}

type ClientGRPC struct {
	Connect         string `yaml:"connect"`
	NegotiationType string `yaml:"negotiation_type" env-default:"plaintext"`
	Cert            string `yaml:"cert"`
	MaxMsgSize      int    `yaml:"max_msg_size" env-default:"4"`
}

type Certs struct {
	Use bool   `yaml:"use"`
	Crt string `yaml:"crt"`
	Key string `yaml:"key"`
}

type SSO struct {
	ClientID      string `yaml:"client_id"`
	ClientSecret  string `yaml:"client_secret"`
	AuthServerUrl string `yaml:"auth_server_url"`
}

var cfg *Config

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if len(configPath) == 0 {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file \"" + configPath + "\" does not exist")
	}

	cfg = new(Config)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic("cant load config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
