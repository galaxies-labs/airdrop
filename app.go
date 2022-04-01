package main

import (
	"fmt"
	"os"

	"github.com/galaxies-labs/airdrop/cmd"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Errorf("Can not load .env"))
	}

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
