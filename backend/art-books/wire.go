package main

import (
	"github.com/joaoprofile/gofi-sdk-go/iam/provider/jwt"

	"github.com/Arthur-155/art-books/domain/contas/adapter"
	"github.com/Arthur-155/art-books/domain/contas/mailer"
	"github.com/Arthur-155/art-books/domain/contas/repository"
	"github.com/Arthur-155/art-books/domain/contas/service"
)

// buildContasService monta repository → service; o repository não toca o banco no construtor.
func buildContasService(cfg Config) (service.ContasService, *jwt.Provider, error) {
	provider, err := jwt.NewProvider(jwt.Config{
		Algorithm: jwt.HS256, Secret: cfg.JWTSecret, AccessTokenTTL: cfg.AccessTokenTTL, Issuer: cfg.JWTIssuer,
	})
	if err != nil {
		return nil, nil, err
	}
	contas := service.NewContasService(
		repository.NewContasRepository(),
		mailer.NewResendMailer(mailer.Config{APIKey: cfg.ResendAPIKey, From: cfg.ResendFromEmail}),
		adapter.NewJWTAccessTokenIssuer(provider, cfg.AccessTokenTTL),
		service.Options{BaseURL: cfg.AppBaseURL, RefreshTokenTTL: cfg.RefreshTokenTTL},
	)
	return contas, provider, nil
}

// buildOperatorService monta só o que o comando operacional usa: o repository.
func buildOperatorService() service.ContasService {
	return service.NewContasService(repository.NewContasRepository(), nil, nil, service.Options{})
}
