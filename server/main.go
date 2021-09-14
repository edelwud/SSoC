package main

import (
	"fmt"
	"server/components/command"
)

func main() {
	_, _ = command.NewCommand("ECHO ping")
	cfg := LoadConfig()
	fmt.Print(cfg.Server.Host)
}
