package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Message struct {
	Handler string
}

type conf struct {
	Test     string
	Messages map[string]Message
}

const File string = "./properties.conf"

var Config conf

func init() {
	_, err := toml.DecodeFile(File, &Config)
	if err != nil {
		fmt.Println(err)
	}
}

func GetDataHandler(target string) string {
	if val, ok := Config.Messages[target]; ok {
		return val.Handler
	}
	return Config.Messages["default"].Handler
}
