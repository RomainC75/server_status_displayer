package settings

var Settings *YAML

type SettingsInterface interface {
	Get() YAML
}

func (s *YAML) Get() YAML {
	return *s
}
