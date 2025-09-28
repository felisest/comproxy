package config

type Config struct {
	RemoteHost  string
	TestingHost string
	Path        string
	Port        string
	Rate        int64
}

func NewConfig() *Config {
	return &Config{
		RemoteHost:  "https://jsonplaceholder.typicode.com",
		TestingHost: "https://jsonplaceholder.typicode.com",
		Path:        "/users",
		Port:        "8082",
		Rate:        1,
	}
}
