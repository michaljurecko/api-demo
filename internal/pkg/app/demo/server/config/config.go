package config

type Config struct {
	ListenAddress string `default:"0.0.0.0:8000" help:"Listen address of API HTTP server."`
}
