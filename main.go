package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"goxenith/app/cmd"
	_ "goxenith/app/models/ent/runtime"
	"goxenith/bootstrap"
	btsConig "goxenith/config"
	"goxenith/pkg/config"
	"goxenith/pkg/console"
	"os"
)

func init() {
	btsConig.Initialize()
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "Goxenith",
		Short: "A community project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,
		PersistentPreRun: func(command *cobra.Command, args []string) {
			config.InitConfig(cmd.Env)
			bootstrap.SetupLogger()
			bootstrap.SetupDB()
			bootstrap.SetupRedis()
		},
	}

	rootCmd.AddCommand(
		cmd.CmdServe,
	)

	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)
	cmd.RegisterGlobalFlags(rootCmd)
	rootCmd.SetHelpCommand(cmd.CobraHelpCommand(rootCmd))
	rootCmd.SetUsageTemplate(cmd.CobraCommandUsageTemplate())
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
