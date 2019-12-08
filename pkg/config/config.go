package config

import (
	"context"

	"github.com/koding/multiconfig"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/constants"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
)

type Config struct {
	Log     LogConfig
	Grpc    GrpcConfig
	Prisma  PrismaConfig
	Etcd    EtcdConfig
	IAM     IAMConfig
	Process ProcessConfig
}

type LogConfig struct {
	Level string `default:"info"` // debug, info, warn, error, fatal
}

type GrpcConfig struct {
	ShowErrorCause bool `default:"false"` // show grpc error cause to frontend
}

type PrismaConfig struct {
	MysqlEndpoint string `default:"http://ltk-mysql-prisma:4466/ltk/mysql"`
	Disable       bool   `default:"false"`
}

type EtcdConfig struct {
	Endpoints string `default:"ltk-etcd:2379"`
}

type IAMConfig struct {
	SecretKey string `default:"test"`
}

type ProcessConfig struct {
}

func LoadConfig() *Config {
	ctx := context.Background()
	config := new(Config)
	m := multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader(constants.EnvLoaderPrefix))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		logger.Critical(ctx, "Failed to load config: %+v", err)
		panic(err)
	}
	logger.SetLevelByString(config.Log.Level)
	logger.Debug(ctx, "LoadConfig: %+v", config)

	return config
}
