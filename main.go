/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

// Install go-bindata if not present. This is a convenience directive for development environments.
//go:generate go install github.com/go-bindata/go-bindata/...
//go:generate go-bindata -o ./internal/config/config.json.go -pkg config ./data/config.json

import (
	"moba-converter-go/cmd"
	_ "moba-converter-go/cmd/config"
	_ "moba-converter-go/cmd/convert"
)

func main() {
	cmd.Execute()
}
