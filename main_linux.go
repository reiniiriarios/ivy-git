//go:build linux

package main

import "github.com/wailsapp/wails/v2/pkg/options"

func appBgColor() *options.RGBA {
	return &options.RGBA{R: 7, G: 18, B: 34, A: 255}
}
