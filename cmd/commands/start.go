package cmd

import (
	"log"

	_interface "github.com/polpettone/streamdeckd/cmd/interface"
	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start Streamdeck",
		Run: func(cmd *cobra.Command, args []string) {
			handleStartCommand(cmd, args)
		},
	}
}

func handleStartCommand(cobraCommand *cobra.Command, args []string) error {

	configPath, err := cobraCommand.Flags().GetString("config")
	if err != nil {
		return err
	}

	deprecatedConfig, err := cobraCommand.Flags().GetBool("deprecated")
	if err != nil {
		return err
	}
	log.Printf("Using deprecatedConfig: %t", deprecatedConfig)

	engine := _interface.NewEngine(configPath)

	return engine.Run(!deprecatedConfig)
}

func init() {
	startCmd := StartCmd()
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP(
		"config",
		"c",
		"",
		"path to config file")

	startCmd.Flags().BoolP(
		"deprecated",
		"d",
		false,
		"use deprecated config")

	err := startCmd.MarkFlagRequired("config")

	if err != nil {
		log.Fatal(err)
	}
}
