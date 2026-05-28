package main

import (
	"fmt"
	"os"
	"strings"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/viper"
)

const (
	authPassword  = "password"
	authSignature = "signature"
)

func newAPIClient() (*regru.Client, error) {
	username := viper.GetString("username")
	password := viper.GetString("password")
	privateKeyPath := viper.GetString("private_key")
	certFile := viper.GetString("cert_file")
	keyFile := viper.GetString("key_file")
	baseURL := viper.GetString("base_url")
	authMode := strings.ToLower(strings.TrimSpace(viper.GetString("auth")))

	if username == "" {
		return nil, fmt.Errorf("укажите REGRU_USERNAME (флаг --username, .env или переменную окружения)")
	}

	if authMode == "" {
		if privateKeyPath != "" {
			authMode = authSignature
		} else {
			authMode = authPassword
		}
	}

	opts := []regru.ClientOption{
		regru.WithBaseURL(baseURL),
		regru.WithOutputFormat(regru.OutputJSON),
	}

	if certFile != "" {
		if keyFile == "" {
			keyFile = privateKeyPath
		}
		if keyFile == "" {
			return nil, fmt.Errorf("для TLS укажите --key-file или --private-key (путь к PEM-ключу)")
		}
		tlsOpt, err := regru.NewTLSClientAuthOption(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		opts = append(opts, tlsOpt)
	}

	switch authMode {
	case authPassword:
		if password == "" {
			return nil, fmt.Errorf("режим password: укажите REGRU_PASSWORD")
		}
		return regru.NewClient(username, password, opts...), nil

	case authSignature:
		if privateKeyPath == "" {
			return nil, fmt.Errorf("режим signature: укажите --private-key или REGRU_PRIVATE_KEY")
		}
		keyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("чтение private key: %w", err)
		}
		return regru.ClientWithSignature(username, keyData, opts...)

	default:
		return nil, fmt.Errorf("неизвестный --auth %q (допустимо: password, signature)", authMode)
	}
}

func authModeLabel() string {
	if a := viper.GetString("auth"); a != "" {
		return a
	}
	if viper.GetString("private_key") != "" {
		return "signature (auto)"
	}
	return "password (auto)"
}
