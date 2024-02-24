package settings

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type YAML struct {
	URLs       []Url `yaml:"URLs"`
	Interval_s int   `yaml:"Interval_s"`
}

type Url struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

func SetSettings() error {
	yfile, err := ioutil.ReadFile("settings.yaml")

	if err != nil {
		return err
	}

	var data YAML

	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		return err
	}
	Settings = &data
	return nil
}
