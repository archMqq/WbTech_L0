package main

import (
	"L0/internal/api"
	"L0/internal/config"
	"L0/internal/database"
)

func main() {
	cfg := config.Load()
	rep := database.Init(cfg)
	api.Init(cfg, rep)
}
