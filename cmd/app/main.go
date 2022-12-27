package main

import (
	"os"

	"github.com/wgo-admin/backend/internal/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd := app.NewAppCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
