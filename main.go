/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"moba-converter-go/cmd"
	_ "moba-converter-go/cmd/config"
	_ "moba-converter-go/cmd/convert"
)

func main() {
	cmd.Execute()
}
