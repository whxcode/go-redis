package config

import (
	"encoding/json"
	"io"
	"os"
)

type ServerProperties struct {
	Bind           string `cfg:bind`
	Port           int    `cfg:port`
	AppendOnly     bool   `cfg:appendOnly`
	AppendFilename string `cfg:appendFilename`
	MaxClients     int    `cfg:maxClients`
	RequirePass    string `cfg:requirePass`
	DataBases      int    `cfg:databases`

	Peers          []string `cfg:peers`
	Self           string   `cfg:self`
	Appendonly     string   `cfg:appendonly`
	appendFilename string   `cfg:appendfilename`
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

func SetupConfig(filename string) {
	// 读取配置文件
	byte, _ := os.ReadFile(filename)

	err := json.Unmarshal(byte, Properties)
	if err != nil {
		panic(err)
	}
}

func MakeServerProperties() *ServerProperties {
	return &ServerProperties{
		Bind:           "0.0.0.0",
		Port:           6378,
		DataBases:      16,
		Appendonly:     "yes",
		AppendFilename: "appendonly.aof",
	}
}

func init() {
	Properties = MakeServerProperties()
}
