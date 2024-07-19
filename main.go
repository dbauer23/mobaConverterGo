/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

//go:generate go-bindata -o .\internal\config\config.json.go -pkg config .\data\config.json

import (
	"moba-converter-go/cmd"
	_ "moba-converter-go/cmd/config"
	_ "moba-converter-go/cmd/convert"
)

func main() {
	cmd.Execute()
}
