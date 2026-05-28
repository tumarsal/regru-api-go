package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Показать баланс аккаунта",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}

		balance, err := client.User().GetBalance(context.Background())
		if err != nil {
			return err
		}

		fmt.Printf("%.2f %s\n", balance.Balance, balance.Currency)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
