package main

import (
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/config"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pi"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/service/process"
)

func main() {
	cfg := config.LoadConfig()
	pi.SetGlobal(cfg)

	process.Serve()
}
