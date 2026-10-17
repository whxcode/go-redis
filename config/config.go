package config

import "io"

type ServerProperties struct {
	Bind           string `cfg:bind`
	Port           int    `cfg:port`
	AppendOnly     bool   `cfg:appendOnly`
	AppendFilename string `cfg:appendFilename`
	MaxClients     int    `cfg:maxClients`
	RequirePass    string `cfg:requirePass`
	DataBases      int    `cfg:databases`

	Peers []string `cfg:peers`
	Self  string   `cfg:self`
}

var Properties *ServerProperties

func init() {
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false,
	}
}

func parse(src io.Reader) *ServerProperties {
	config := &ServerProperties{}

	return config
}
