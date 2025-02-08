package main

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/cmd"
	"github.com/NesterovYehor/txtnest-cli/config"
)

func main() {
	_ = config.Init()
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
