package handlers

import (
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations/roles"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Role struct {
	client *ent.Client
	logger *zap.Logger
}

func NewRole(client *ent.Client, logger *zap.Logger) *Role {
	return &Role{
		client: client,
		logger: logger,
	}
}

func (r Role) GetRolesFunc() roles.GetRolesHandlerFunc {
	return func(s roles.GetRolesParams) middleware.Responder {
		e, err := r.client.Role.Query().Order(ent.Asc("id")).All(s.HTTPRequest.Context())
		if err != nil {
			return roles.NewGetRolesDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
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