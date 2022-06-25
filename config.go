package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ModeratorChat int64
	ChannelChat   int64
	Token         string
}

func (config *Config) getConfig() (err error) {
	d, e := ioutil.ReadFile("config.json")
	if e != nil {
		config = nil
		return e
	}
	e = json.Unmarshal(d, &config)
	return e
}
