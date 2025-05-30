package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/CSR-LC/csr-be/internal/generated/swagger/models"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations/password_reset"
	"github.com/CSR-LC/csr-be/internal/messages"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

func SetPasswordResetHandler(logger *zap.Logger, api *operations.BeAPI, passwordService domain.PasswordResetService) {
	PasswordResetHandler := NewPasswordReset(logger, passwordService)

	api.PasswordResetSendLinkByLoginHandler = PasswordResetHandler.SendLinkByLoginFunc()
	api.PasswordResetGetPasswordResetLinkHandler = PasswordResetHandler.GetPasswordResetLinkFunc()
}

type passwordResetHandler struct {
	logger        *zap.Logger
	passwordReset domain.PasswordResetService
}

func NewPasswordReset(logger *zap.Logger, passwordService domain.PasswordResetService) *passwordResetHandler {
	return &passwordResetHandler{
		logger:        logger,
		passwordReset: passwordService,
	}
}

func (c passwordResetHandler) SendLinkByLoginFunc() password_reset.SendLinkByLoginHandlerFunc {
	return func(s password_reset.SendLinkByLoginParams) middleware.Responder {
		ctx := s.HTTPRequest.Context()
		login := *s.Login.Data.Login
		if login == "" {
			c.logger.Warn("Login is empty")
			return password_reset.NewSendLinkByLoginDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload(messages.ErrLoginRequired, ""))
		}
		err := c.passwordReset.SendResetPasswordLink(ctx, login)
		if err != nil {
			c.logger.Error("Error while sending reset password link", zap.Error(err))
			return password_reset.NewSendLinkByLoginOK().WithPayload(
				models.PasswordResetResponse(messages.MsgPasswordResetSuccesful))
		}
		return password_reset.NewSendLinkByLoginOK().WithPayload(
			models.PasswordResetResponse(messages.MsgPasswordResetSuccesful))
	}
}

func (c passwordResetHandler) GetPasswordResetLinkFunc() password_reset.GetPasswordResetLinkHandlerFunc {
	return func(s password_reset.GetPasswordResetLinkParams) middleware.Responder {
		ctx := s.HTTPRequest.Context()
		token := s.Token
		err := c.passwordReset.VerifyTokenAndSendPassword(ctx, token)
		if err != nil {
			c.logger.Error("Failed to verify token or send email", zap.Error(err))
			return password_reset.NewGetPasswordResetLinkOK().WithPayload("Check your email for a new password")
		}
		return password_reset.NewGetPasswordResetLinkOK().WithPayload("Check your email for a new password")
	}
}
