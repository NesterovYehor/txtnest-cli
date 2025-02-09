package main

import (
	"database/sql"
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/cmd"
	"github.com/NesterovYehor/txtnest-cli/config"
	"github.com/NesterovYehor/txtnest-cli/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	tokensDB, err := sql.Open("sqlite3", "./tokens.db")
	if err != nil {
		fmt.Printf("Failed to conect to sqlite: %v", err)
		return
	}

	if err := storage.InitTokenStorage(tokensDB); err != nil {
		fmt.Println(err)
		return
	}
	_ = config.Init()
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
