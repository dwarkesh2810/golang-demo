package conf

import "log"

var Env Config

func LoadEnv(path string) error {
	EnvConfig, err := LoadConfig(path, "app", "env")
	if err != nil {
		log.Fatal("failed to change ", err)
		return err
	}
	Env = EnvConfig
	return nil
}
