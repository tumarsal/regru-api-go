package main

import (
	"context"
	"fmt"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/cobra"
)

var nopCmd = &cobra.Command{
	Use:   "nop",
	Short: "Проверить доступность API",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		testClient := regru.NewClient("test", "test")
		if err := testClient.Nop(ctx); err != nil {
			return fmt.Errorf("test nop: %w", err)
		}
		fmt.Println("API доступен (test/test)")

		client, err := newAPIClient()
		if err != nil {
			return err
		}
		if err := client.Nop(ctx); err != nil {
			return fmt.Errorf("authenticated nop: %w", err)
		}
		fmt.Println("API доступен (ваши учётные данные)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nopCmd)
}
