package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/CSR-LC/csr-be/internal/generated/swagger/models"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations"
	"github.com/CSR-LC/csr-be/internal/generated/swagger/restapi/operations/roles"
	"github.com/CSR-LC/csr-be/internal/messages"
	"github.com/CSR-LC/csr-be/internal/repositories"
	"github.com/CSR-LC/csr-be/pkg/domain"
)

func SetRoleHandler(logger *zap.Logger, api *operations.BeAPI) {
	roleRepo := repositories.NewRoleRepository()
	roleHandler := NewRole(logger)

	api.RolesGetRolesHandler = roleHandler.GetRolesFunc(roleRepo)
}

type Role struct {
	logger *zap.Logger
}

func NewRole(logger *zap.Logger) *Role {
	return &Role{
		logger: logger,
	}
}

func (r Role) GetRolesFunc(repository domain.RoleRepository) roles.GetRolesHandlerFunc {
	return func(s roles.GetRolesParams, _ *models.Principal) middleware.Responder {
		ctx := s.HTTPRequest.Context()
		e, err := repository.GetRoles(ctx)
		if err != nil {
			r.logger.Error(messages.ErrQueryRoles)
			return roles.NewGetRolesDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(messages.ErrQueryRoles, ""))
		}
		listRoles := models.ListRoles{}
		for _, element := range e {
			id := int64(element.ID)
			listRoles = append(listRoles, &models.Role{
				ID:   &id,
				Name: &element.Name,
				Slug: &element.Slug,
			})
		}
		return roles.NewGetRolesOK().WithPayload(listRoles)
	}
}
