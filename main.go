package main

import (
	"database/sql"
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/cmd"
	"github.com/NesterovYehor/txtnest-cli/config"
	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/NesterovYehor/txtnest-cli/internal/models"
	"github.com/NesterovYehor/txtnest-cli/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

func init() {
	_ = config.Init()
	tokensDB, err := sql.Open("sqlite3", "./tokens.db")
	if err != nil {
		fmt.Printf("Failed to conect to sqlite: %v", err)
		return
	}

	if err := storage.InitTokenStorage(tokensDB); err != nil {
		fmt.Println(err)
		return
	}
	client := api.GetInstance()
	tokenStore, err := storage.GetTokenStorage()
	if err != nil {
		fmt.Println(err)
	}
	client.RegisterTokenUpdateCallback(func(newTokens *models.TokenData) error {
		return tokenStore.SaveTokens(newTokens)
	})
}
