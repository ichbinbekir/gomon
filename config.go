package gomon

import "reflect"

type Config struct {
	// TODO Comment: Events and Errors channel buffer size, both fsnotify and gomon.
	// Events channel buffer size, both fsnotify and gomon.
	BufferSize uint

	Save string
}

var defaultConfig = Config{
	BufferSize: 50,
}

func mergeConfigs(cfgs ...Config) Config {
	config := defaultConfig
	if len(cfgs) == 0 {
		return config
	}

	configValue := reflect.ValueOf(&config).Elem()
	nField := configValue.NumField()
	for _, cfg := range cfgs {
		cfgValue := reflect.ValueOf(cfg)
		for field := range nField {
			value := cfgValue.Field(field)
			if !value.IsZero() {
				configValue.Field(field).Set(value)
			}
		}
	}

	return config
}
