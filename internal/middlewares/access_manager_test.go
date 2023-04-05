package middlewares

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/authentication"
)

const (
	userRole                  = "user"
	adminRole                 = "admin"
	simpleValidPath           = "/v1/simple"
	validPathWithParam        = "/v1/with/params/{id}"
	validPathWithParamExample = "/v1/with/params/1"
	validPathWithParams       = "/v1/with/params/{id}/and/{name}"
	simpleInvalidPath         = "/v1/invalid"
)

func Test_blackListAccessManager(t *testing.T) {
	var manager AccessManager
	t.Run("NewAccessManager", func(t *testing.T) {
		roles := []Role{
			{
				Slug: userRole,
			},
			{
				Slug: adminRole,
			},
		}
		fullAccessRoles := []Role{
			{
				Slug: adminRole,
			},
		}
		endpoints := ExistingEndpoints{
			http.MethodGet: {
				simpleValidPath,
				validPathWithParam,
				validPathWithParams,
			},
		}
		logger := zap.NewNop()
		var err error
		manager, err = NewAccessManager(roles, fullAccessRoles, endpoints, logger)
		assert.NoError(t, err)
	})

	t.Run("AddNewAccess", func(t *testing.T) {
		type accessRule struct {
			role   Role
			method string
			path   string
			isErr  bool
			isOk   bool
		}
		newAccessRules := []accessRule{
			{
				role:   Role{Slug: userRole, IsEmailConfirmed: true},
				method: http.MethodGet,
				path:   simpleValidPath,
				isOk:   true,
			},
			{
				role:   Role{Slug: userRole, IsEmailConfirmed: true},
				method: http.MethodGet,
				path:   simpleValidPath + "/",
				isErr:  true,
			},
			{
				role:   Role{Slug: userRole, IsEmailConfirmed: true},
				method: http.MethodGet,
				path:   simpleInvalidPath,
				isErr:  true,
			},
			{
				role:   Role{Slug: userRole, IsEmailConfirmed: true},
				method: http.MethodGet,
				path:   validPathWithParam,
				isOk:   true,
			},
			{
				role:   Role{Slug: userRole, IsEmailConfirmed: true},
				method: http.MethodPut,
				path:   validPathWithParam,
				isErr:  true,
			},
			{
				role:   Role{Slug: adminRole},
				method: http.MethodGet,
				path:   validPathWithParams,
				isOk:   false,
				isErr:  false,
			},
			{
				role:   Role{Slug: userRole, IsEmailConfirmed: true},
				method: http.MethodGet,
				path:   simpleInvalidPath,
				isErr:  true,
			},
			{
				role:   Role{Slug: "unknown", IsEmailConfirmed: true},
				method: http.MethodGet,
				path:   validPathWithParams,
				isErr:  true,
			},
		}
		for _, rule := range newAccessRules {
			ok, err := manager.AddNewAccess(rule.role, rule.method, rule.path)
			if rule.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equalf(t, rule.isOk, ok, "AddNewAccess(%v, %s, %s)", rule.role, rule.method, rule.path)
		}
	})

	type requestData struct {
		role         Role
		method, path string
		hasAccess    bool
	}
	requestsData := []requestData{
		{
			role:      Role{Slug: userRole, IsEmailConfirmed: true},
			method:    http.MethodGet,
			path:      endpointConversion(simpleValidPath),
			hasAccess: true,
		},
		{
			role:      Role{Slug: adminRole, IsEmailConfirmed: true},
			method:    http.MethodGet,
			path:      endpointConversion(simpleValidPath),
			hasAccess: true,
		},
		{
			role:      Role{Slug: userRole, IsEmailConfirmed: true},
			method:    http.MethodGet,
			path:      endpointConversion(validPathWithParamExample),
			hasAccess: true,
		},
		{
			role:      Role{Slug: userRole, IsEmailConfirmed: true},
			method:    http.MethodGet,
			path:      strings.TrimPrefix(endpointConversion(validPathWithParamExample), "/"),
			hasAccess: true,
		},
		{
			role:      Role{Slug: userRole, IsEmailConfirmed: true},
			method:    http.MethodGet,
			path:      strings.TrimPrefix(endpointConversion(validPathWithParamExample), "/") + "/",
			hasAccess: true,
		},
		{
			role:      Role{Slug: userRole, IsEmailConfirmed: true},
			method:    http.MethodPut,
			path:      endpointConversion(validPathWithParamExample),
			hasAccess: false,
		},
	}

	t.Run("HasAccess", func(t *testing.T) {
		for _, data := range requestsData {
			assert.Equalf(t, data.hasAccess, manager.HasAccess(data.role, data.method, data.path),
				"HasAccess(%v, %s, %s)", data.role, data.method, data.path)
		}
	})

	t.Run("Authorize", func(t *testing.T) {
		for _, data := range requestsData {
			request := &http.Request{
				Method: data.method,
				URL: &url.URL{
					Path: data.path,
				},
			}
			auth := authentication.Auth{
				Role: &authentication.Role{
					Slug: data.role.Slug,
				},
				IsEmailConfirmed: true,
			}
			err := manager.Authorize(request, auth)
			if data.hasAccess {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		}

	})
}
