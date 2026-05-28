package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [domain]",
	Short: "Проверить доступность домена",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		isPremium, _ := cmd.Flags().GetBool("premium")
		result, err := client.Domain().Check(context.Background(), args[0], isPremium)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", result.DomainName, result.Result)
		if result.IsPremium && result.Price > 0 {
			fmt.Printf("premium price: %.2f\n", result.Price)
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().Bool("premium", false, "запросить цену премиум-домена")
	rootCmd.AddCommand(checkCmd)
}
