package config

import "time"

type Configuration struct {
	Proxy Proxy
}

type Proxy struct {
	RemoteHost  string
	TestingHost string
	Path        string
	Rate        int64 //Trigger rate 1:N
	Server      Server
}

type Server struct {
	Port    string
	Timeout int
}

func (w *Server) GetTimeout() time.Duration {
	return time.Duration(w.Timeout) * time.Second
}
