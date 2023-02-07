package main

import (
	"fmt"
	"os"

	"gpe_project/internal/app/adapter"
	"gpe_project/internal/app/adapter/service"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("fatal error during dir finding: %w", err))
	}

	service.LoadConfig(dir)
	service.NewLogging()

	engine := adapter.Setup()
	if err := engine.Run(); err != nil {
		panic(fmt.Errorf("fatal error Gin server: %w", err))
	}
}
