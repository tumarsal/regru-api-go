package main

import "github.com/spf13/cobra"

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Управление DNS-зоной домена",
}

func init() {
	rootCmd.AddCommand(zoneCmd)
}
