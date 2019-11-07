package main

// Config holds TySug's central config parameters
type Config struct {
	References map[string][]string `toml:"references"`
	Client     struct {
		InputLengthMax int `toml:"inputLengthMax"`
	} `toml:"client"`
	CORS struct {
		AllowedOrigins []string `toml:"allowedOrigins"`
	} `toml:"CORS"`
	Server struct {
		ListenOn string `toml:"listenOn"`
		Headers  []struct {
			Name  string `toml:"name"`
			Value string `toml:"value"`
		} `toml:"headers"`
		Log struct {
			Level string `toml:"level"`
		} `toml:"log"`
		Profiler struct {
			Enable bool   `toml:"enable"`
			Prefix string `toml:"prefix"`
		} `toml:"profiler"`
	} `toml:"server"`
}
