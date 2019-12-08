package main

import (
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/apigateway"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/config"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pi"
)

func main() {
	cfg := config.LoadConfig()
	cfg.Prisma.Disable = true

	pi.SetGlobal(cfg)
	apigateway.Serve()
}
