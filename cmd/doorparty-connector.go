package main

import (
	"github.com/echicken/dpc2/internal/config"
	"github.com/echicken/dpc2/internal/server"
	"github.com/echicken/dpc2/internal/tunnel"
)

func main() {
	cfg := config.Get()
	server.Listen(cfg, tunnel.Start)
}
