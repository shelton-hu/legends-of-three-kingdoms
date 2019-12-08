package config

import (
	"fmt"

	"github.com/shelton-hu/util/yamlutil"
)

type GlobalConfig struct {
	Common  CommonGlobalConfig  `json:"common"`
	IAM     IAMGlobalConfig     `json:"iam"`
	Process ProcessGlobalConfig `json:"process"`
}

// common
type CommonGlobalConfig struct {
}

// iam
type IAMGlobalConfig struct {
}

// sync
type ProcessGlobalConfig struct {
}

func init() {
	DecodeInitConfig()
}

func DecodeInitConfig() GlobalConfig {
	globalConfig, err := ParseGlobalConfig([]byte(InitialGlobalConfig))
	if err != nil {
		fmt.Print("InitialGlobalConfig is invalid, please fix it")
		panic(err)
	}
	return globalConfig
}

func EncodeGlobalConfig(conf GlobalConfig) string {
	out, err := yamlutil.Encode(conf)
	if err != nil {
		fmt.Print("Encode globalConfig failed")
		panic(err)
	}
	return string(out)
}

func ParseGlobalConfig(data []byte) (GlobalConfig, error) {
	var globalConfig GlobalConfig
	err := yamlutil.Decode(data, &globalConfig)
	if err != nil {
		return globalConfig, err
	}
	return globalConfig, nil
}
