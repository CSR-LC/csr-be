package handlers

import (
	"net/http"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
)

func buildErrorPayload(code int32, msg string) *models.SwaggerError {
	return &models.SwaggerError{
		Code:     &code,
		Message:  &msg,
	}
}

func buildInternalErrorPayload(msg string) *models.SwaggerError {
	return buildErrorPayload(http.StatusInternalServerError, msg)
}

func buildExFailedErrorPayload(msg string) *models.SwaggerError {
	return buildErrorPayload(http.StatusExpectationFailed, msg)
}

func buildConflictErrorPayload(msg string) *models.SwaggerError {
	return buildErrorPayload(http.StatusConflict, msg)
}

func buildNotFoundErrorPayload(msg string) *models.SwaggerError {
	return buildErrorPayload(http.StatusNotFound, msg)
}

func buildForbiddenErrorPayload(msg string) *models.SwaggerError {
	return buildErrorPayload(http.StatusForbidden, msg)
}

func buildBadRequestErrorPayload(msg string) *models.SwaggerError {
	return buildErrorPayload(http.StatusBadRequest, msg)
}
