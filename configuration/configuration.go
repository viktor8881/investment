package configuration

import (
	"encoding/json"
	//"log"
	"os"
)

type Favorit struct {
	FIGI         string
	ExchangeName string
}

type configuration struct {
	Sdk struct {
		Token    string
		Favorits []Favorit
	}
	Email struct {
		Host     string
		Port     int64
		Username string
		Password string
		To       []string
		From     string
	}
	LogFile       string
	HostName      string
	PrefixPathImg string
}

func LoadConfiguration(fileName string) (*configuration, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	configuration := &configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}
