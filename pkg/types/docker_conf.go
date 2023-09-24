package types

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/khulnasoft-lab/fanal/types"
)

type DockerConfig struct {
	UserName string `env:"VUL_USERNAME"`
	Password string `env:"VUL_PASSWORD"`
	Insecure bool   `env:"VUL_INSECURE" envDefault:"true"`

	DockerCertPath string `env:"DOCKER_CERT_PATH"`
	DockerHost     string `env:"DOCKER_HOST"`
}

func GetDockerOption(timeout time.Duration) (types.DockerOption, error) {
	cfg := DockerConfig{}
	if err := env.Parse(&cfg); err != nil {
		return types.DockerOption{}, err
	}

	return types.DockerOption{
		UserName:              cfg.UserName,
		Password:              cfg.Password,
		Timeout:               timeout,
		InsecureSkipTLSVerify: cfg.Insecure,
		DockerDaemonCertPath:  cfg.DockerCertPath,
		DockerDaemonHost:      cfg.DockerHost,
	}, nil
}
