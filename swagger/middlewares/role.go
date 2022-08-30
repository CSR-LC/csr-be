package middlewares

import (
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/authentication"
)

const roleField = "role"

type roleClaim struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

func CheckRole(target authentication.Slug, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if len(auth) == 0 {
				logger.Warn("HTTP request doesn't contain Authorization header", zap.Any("headers", r.Header))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h := strings.SplitN(auth, " ", 2)
			if len(h) != 2 {
				logger.Warn("Invalid Authorization header", zap.Any("auth_header", auth))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if strings.ToLower(h[0]) != "bearer" {
				logger.Warn("Invalid authorization schema", zap.Any("schema", h[0]))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			jwt.RegisterCustomField(roleField, roleClaim{})
			token, err := jwt.ParseString(h[0])
			if err != nil {
				logger.Warn("Failed to parse JWT token", zap.Any("token", token))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			v, ok := token.Get(roleField)
			if !ok {
				logger.Warn("Failed to get role from token")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			role, ok := v.(roleClaim)
			if !ok {
				logger.Warn("Failed to extract the role")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if role.Slug != string(target) {
				logger.Warn("Not enought permissions")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
