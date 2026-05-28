package main

import (
	"fmt"
	"os"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/viper"
)

func newAPIClient() (*regru.Client, error) {
	username := viper.GetString("username")
	password := viper.GetString("password")
	privateKeyPath := viper.GetString("private_key")
	baseURL := viper.GetString("base_url")

	if username == "" {
		return nil, fmt.Errorf("укажите REGRU_USERNAME (флаг --username, .env или переменную окружения)")
	}

	opts := []regru.ClientOption{
		regru.WithBaseURL(baseURL),
		regru.WithOutputFormat(regru.OutputJSON),
	}

	if privateKeyPath != "" {
		keyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("чтение private key: %w", err)
		}
		return regru.ClientWithSignature(username, keyData, opts...)
	}

	if password == "" {
		return nil, fmt.Errorf("укажите REGRU_PASSWORD (флаг --password, .env или переменную окружения)")
	}

	return regru.NewClient(username, password, opts...), nil
}
