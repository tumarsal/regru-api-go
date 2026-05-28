package main

import (
	"fmt"
	"os"
	"strings"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "regru",
		Short: "CLI для REG.RU API 2.0",
		Long:  "Клиент командной строки для работы с REG.RU API (домены, баланс, DNS и др.).",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "путь к .env (по умолчанию: .env в текущей директории)")
	rootCmd.PersistentFlags().String("auth", "",
		"способ авторизации: password (логин+пароль) или signature (RSA-подпись); по умолчанию: password, либо signature если задан --private-key")
	rootCmd.PersistentFlags().String("username", "", "логин REG.RU (env: REGRU_USERNAME)")
	rootCmd.PersistentFlags().String("password", "", "пароль REG.RU (env: REGRU_PASSWORD)")
	rootCmd.PersistentFlags().String("private-key", "", "путь к PEM-ключу RSA для подписи (env: REGRU_PRIVATE_KEY)")
	rootCmd.PersistentFlags().String("cert-file", "", "клиентский TLS-сертификат .crt/.pem (env: REGRU_CERT_FILE)")
	rootCmd.PersistentFlags().String("key-file", "", "ключ TLS-сертификата; по умолчанию --private-key (env: REGRU_KEY_FILE)")
	rootCmd.PersistentFlags().String("base-url", regru.DefaultBaseURL,
		fmt.Sprintf("базовый URL API (по умолчанию: %s)", regru.DefaultBaseURL))

	viper.SetDefault("base_url", regru.DefaultBaseURL)
}

func initConfig() {
	viper.SetEnvPrefix("REGRU")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	_ = viper.BindEnv("auth", "REGRU_AUTH")
	_ = viper.BindEnv("username", "REGRU_USERNAME")
	_ = viper.BindEnv("password", "REGRU_PASSWORD")
	_ = viper.BindEnv("private_key", "REGRU_PRIVATE_KEY")
	_ = viper.BindEnv("cert_file", "REGRU_CERT_FILE")
	_ = viper.BindEnv("key_file", "REGRU_KEY_FILE")
	_ = viper.BindEnv("base_url", "REGRU_BASE_URL")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if cfgFile != "" {
			fmt.Fprintf(os.Stderr, "не удалось прочитать конфиг %q: %v\n", cfgFile, err)
			os.Exit(1)
		}
		// .env не обязателен — можно использовать REGRU_* или флаги
	}

	normalizeConfigKeys()
	bindChangedFlags()
}

// normalizeConfigKeys сопоставляет REGRU_* из .env с ключами viper (username, password, ...).
func normalizeConfigKeys() {
	aliases := map[string]string{
		"regru_auth":        "auth",
		"regru_username":    "username",
		"regru_password":    "password",
		"regru_private_key": "private_key",
		"regru_cert_file":   "cert_file",
		"regru_key_file":    "key_file",
		"regru_base_url":    "base_url",
	}
	for src, dst := range aliases {
		if viper.GetString(dst) != "" {
			continue
		}
		if v := viper.GetString(src); v != "" {
			viper.Set(dst, v)
		}
	}
}

// bindChangedFlags применяет только явно переданные флаги (пустые не затирают .env).
func bindChangedFlags() {
	flagToKey := map[string]string{
		"auth":        "auth",
		"username":    "username",
		"password":    "password",
		"private-key": "private_key",
		"cert-file":   "cert_file",
		"key-file":    "key_file",
		"base-url":    "base_url",
	}
	rootCmd.PersistentFlags().Visit(func(f *pflag.Flag) {
		key, ok := flagToKey[f.Name]
		if ok && f.Changed {
			viper.Set(key, f.Value.String())
		}
	})
}
