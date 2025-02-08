package main

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

