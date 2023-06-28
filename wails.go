package main

import (
	_ "embed"
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed wails.json
var wailsJson []byte

type wailsConfig struct {
	Name   string            `json:"name"`
	Slug   string            `json:"slug"`
	Bundle string            `json:"bundle"`
	Author wailsConfigAuthor `json:"author"`
	Info   wailsConfigInfo   `json:"info"`
}

type wailsConfigAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type wailsConfigInfo struct {
	Version string `json:"productVersion"`
}

func (a *App) getWailsConfig() wailsConfig {
	var cfg wailsConfig
	if err := json.Unmarshal(wailsJson, &cfg); err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	return cfg
}
