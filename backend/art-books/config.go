package main

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/joaoprofile/gofi-sdk-go/base/common"
	"github.com/joaoprofile/gofi-sdk-go/base/environment"
	"github.com/joaoprofile/gofi-sdk-go/base/secrets"
)

const (
	minJWTSecretBytes = 32
	// trustPrivateNetworks é a palavra do SDK que confia nas faixas privadas.
	trustPrivateNetworks = "private"
)

// Config reúne o que o serviço lê do ambiente; só este arquivo toca o processo.
type Config struct {
	AllowedOrigins  []string
	AppBaseURL      string `env:"APP_BASE_URL"`
	ResendAPIKey    string `env:"RESEND_API_KEY"`
	ResendFromEmail string `env:"RESEND_FROM_EMAIL"`
	// TrustedProxies decide de onde o X-Forwarded-For vale: o bloqueio por IP depende dele.
	TrustedProxies string `env:"HTTP_TRUSTED_PROXIES"`

	JWTSecret      []byte
	JWTIssuer       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// LoadConfig devolve todos os erros juntos; as três variáveis do Resend e do link são obrigatórias.
func LoadConfig(ctx context.Context, env *environment.Environment) (Config, error) {
	auth := env.Auth()
	cfg := Config{
		AllowedOrigins: env.HTTP().AllowedOrigins,
		JWTSecret:      auth.JWTSecret, JWTIssuer: auth.Issuer,
		AccessTokenTTL: auth.AccessTokenTTL, RefreshTokenTTL: auth.RefreshTokenTTL,
	}
	resolver := secrets.NewResolver()
	var failures []error
	lookup := func(key string) string {
		value, err := resolver.Resolve(ctx, os.Getenv(key))
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", key, err))
		}
		return value
	}
	if err := common.ParseStructAnnotationFunc(&cfg, "env", lookup); err != nil {
		failures = append(failures, err)
	}
	failures = append(failures, validateConfig(cfg), environment.LoadError())
	return cfg, errors.Join(failures...)
}

func validateConfig(cfg Config) error {
	var failures []error
	for name, value := range map[string]string{
		"APP_BASE_URL": cfg.AppBaseURL, "RESEND_API_KEY": cfg.ResendAPIKey, "RESEND_FROM_EMAIL": cfg.ResendFromEmail,
	} {
		if value == "" {
			failures = append(failures, fmt.Errorf("%s is required", name))
		}
	}
	if len(cfg.JWTSecret) < minJWTSecretBytes {
		failures = append(failures, fmt.Errorf("JWT_SECRET must have at least %d bytes", minJWTSecretBytes))
	}
	failures = append(failures, validateTrustedProxies(cfg.TrustedProxies))
	return errors.Join(failures...)
}

// validateTrustedProxies aceita CIDRs e "private"; recusa entrada inválida e a que confia em qualquer origem.
func validateTrustedProxies(csv string) error {
	for _, entry := range strings.Split(csv, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" || entry == trustPrivateNetworks {
			continue
		}
		prefix, err := netip.ParsePrefix(entry)
		if err != nil {
			return fmt.Errorf("HTTP_TRUSTED_PROXIES has an invalid entry (%q)", entry)
		}
		if prefix.Bits() == 0 {
			return fmt.Errorf("HTTP_TRUSTED_PROXIES must not trust every address (%s)", prefix)
		}
	}
	return nil
}
