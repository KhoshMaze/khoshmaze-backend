package config

import (
	"encoding/json"
	"os"
)

func readConfig(path string) (Config, error) {
	var c Config

	data, err := os.ReadFile(path)

	if err != nil {
		return c, err
	}

	return c, json.Unmarshal(data, &c)

}

func MustReadConfig(path string) Config {

	conf, err := readConfig(path)

	if err != nil {
		panic(err)
	}

	return conf

}
