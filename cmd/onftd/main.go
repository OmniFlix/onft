package main

import (
	"os"

	"github.com/OmniFlix/onft/app"
	"github.com/OmniFlix/onft/cmd/onftd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "ONFTD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
