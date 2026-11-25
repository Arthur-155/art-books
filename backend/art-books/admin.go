package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/joaoprofile/gofi-sdk-go/base/errs"
	"github.com/joaoprofile/gofi-sdk-go/gofi"
	"github.com/joaoprofile/gofi-sdk-go/gofi/component/database"

	"github.com/Arthur-155/art-books/domain/contas/service"
)

const (
	grantSubcommand  = "grant"
	revokeSubcommand = "revoke"
	adminArgsCount   = 2
)

var errAdminUsage = errors.New("uso: art-books admin grant|revoke <e-mail>")

// runAdmin concede ou retira o papel de administrador pela mesma conexão do serviço (REQ-15).
//
// gofi:implements REQ-15
func runAdmin(ctx context.Context, args []string) error {
	if len(args) != adminArgsCount || (args[0] != grantSubcommand && args[0] != revokeSubcommand) {
		return errAdminUsage
	}
	app, err := gofi.New(serviceName).With(database.New()).Build()
	if err != nil {
		return err
	}
	appErr := changeAdminRole(ctx, buildOperatorService(), args[0], args[1])
	return errors.Join(adminFailure(appErr), app.Shutdown(ctx))
}

func changeAdminRole(ctx context.Context, contas service.ContasService, subcommand, email string) errs.AppError {
	if subcommand == grantSubcommand {
		return contas.GrantAdmin(ctx, email)
	}
	return contas.RevokeAdmin(ctx, email)
}

func adminFailure(appErr errs.AppError) error {
	if !appErr.Exists() {
		return nil
	}
	if errors.Is(appErr, service.ErrContasUserNotFound) {
		return errors.New("conta não encontrada")
	}
	return fmt.Errorf("falha ao alterar o papel: %w", appErr)
}
