package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/user"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/authentication"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	client *ent.Client
	logger *zap.Logger
}

func generateJWT(user *ent.User, jwtSecretKey string, ctx context.Context) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["login"] = user.Login
	claims["role"] = nil
	claims["group"] = nil
	role, err := user.QueryRole().First(ctx)
	if err == nil {
		claims["role"] = map[string]interface{}{
			"id":   role.ID,
			"slug": role.Slug,
		}
	}
	group, err := user.QueryGroups().First(ctx)
	if err == nil {
		claims["group"] = map[string]interface{}{
			"id": group.ID,
		}
	}
	claims["exp"] = time.Now().Add(time.Minute * 300).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func NewUser(client *ent.Client, logger *zap.Logger) *User {
	return &User{
		client: client,
		logger: logger,
	}
}

func buildErrorPayload(err error) *models.Error {
	return &models.Error{
		Data: &models.ErrorData{
			Message: err.Error(),
		},
	}
}

func (c User) LoginUserFunc(jwtSecretKey string) users.LoginHandlerFunc {
	return func(p users.LoginParams) middleware.Responder {
		login := p.Login.Login
		foundUser, err := c.client.User.Query().Where(user.Login(*login)).First(p.HTTPRequest.Context())
		if ent.IsNotFound(err) {
			return users.NewLoginNotFound()
		}
		if err != nil {
			return users.NewLoginDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(*p.Login.Password))
		if err != nil {
			return users.NewLoginNotFound()
		}

		token, err := generateJWT(foundUser, jwtSecretKey, p.HTTPRequest.Context())
		if err != nil {
			return users.NewLoginDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		return users.NewLoginOK().WithPayload(&models.LoginSuccessResponse{
			Data: &models.LoginSuccessResponseData{Token: &token},
		})
	}
}

func (c User) PostUserFunc() users.PostUserHandlerFunc {
	return func(p users.PostUserParams) middleware.Responder {
		// e, err := c.client.User.Create().SetLogin("test").SetEmail("example@example.com").SetPassword("123456").Save(p.HTTPRequest.Context())
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*p.Data.Password), bcrypt.DefaultCost)
		if err != nil {
			return users.NewPostUserDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}
		login := *p.Data.Login
		createdUser, err := c.client.User.
			Create().
			SetEmail(login).
			SetLogin(login).
			SetName(login).
			SetType(user.TypePerson).
			SetPassword(string(hashedPassword)).
			Save(p.HTTPRequest.Context())
		if err != nil {
			if ent.IsConstraintError(err) {
				return users.NewPostUserDefault(http.StatusExpectationFailed).WithPayload(
					buildErrorPayload(errors.New("This login is already used")),
				)
			}
			return users.NewPostUserDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}

		id := int64(createdUser.ID)
		return users.NewPostUserCreated().WithPayload(&models.CreateUserResponse{
			Data: &models.CreateUserResponseData{
				ID: &id,
			},
		})
	}
}

func (c User) UpdateUserByIDFunc() users.UserUpdateHandlerFunc {
	return func(p users.UserUpdateParams) middleware.Responder {
		id := int(p.UserID)
		user, err := c.client.User.Get(p.HTTPRequest.Context(), id)
		if err != nil {
			return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		data := p.UpdateUserTask.Data

		if data.Name != "" {
			user, err = c.client.User.UpdateOne(user).SetName(data.Name).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Surname != "" {
			user, err = c.client.User.UpdateOne(user).SetSurname(data.Surname).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}
		if data.Patronymic != "" {
			user, err = c.client.User.UpdateOne(user).SetPatronymic(data.Patronymic).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.PassportNumber != "" {
			user, err = c.client.User.UpdateOne(user).SetPassportNumber(data.PassportNumber).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.PassportSeries != "" {
			user, err = c.client.User.UpdateOne(user).SetPassportSeries(data.PassportSeries).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.PassportAuthority != "" {
			user, err = c.client.User.UpdateOne(user).SetPassportAuthority(data.PassportAuthority).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.PassportIssueDate != "" {
			t, err := time.Parse("2006-01-02", data.PassportIssueDate)
			user, err = c.client.User.UpdateOne(user).SetPassportIssueDate(t).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Type != string(user.Type) {
			newType := user.Type
			if data.Name == "person" {
				newType = "person"
			} else {
				newType = "organization"
			}
			user, err = c.client.User.UpdateOne(user).SetType(newType).SetName(data.Name).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.PhoneNumber != "" {
			user, err = c.client.User.UpdateOne(user).SetPhone(data.PhoneNumber).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if len(data.ActiveAreas) > 0 {
			var checkMap = make(map[int64]bool, 1)
			newAreas := user.ActiveAreas
			for i := range user.ActiveAreas {
				checkMap[user.ActiveAreas[i]] = true
			}
			for _, v := range data.ActiveAreas {
				if !checkMap[v] {
					checkMap[v] = true
					newAreas = append(newAreas, v)
				}
			}
			user, err = c.client.User.UpdateOne(user).SetActiveAreas(newAreas).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Email != "" {
			user, err = c.client.User.UpdateOne(user).SetEmail(data.Email).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.OrgName != "" {
			user, err = c.client.User.UpdateOne(user).SetOrgName(data.OrgName).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Vk != "" {
			user, err = c.client.User.UpdateOne(user).SetVk(data.Vk).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Instagram != "" {
			user, err = c.client.User.UpdateOne(user).SetInstagram(data.Instagram).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Facebook != "" {
			user, err = c.client.User.UpdateOne(user).SetFacebook(data.Facebook).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}

		if data.Tiktok != "" {
			user, err = c.client.User.UpdateOne(user).SetTiktok(data.Tiktok).Save(p.HTTPRequest.Context())
			if err != nil {
				return users.NewUserUpdateDefault(http.StatusInternalServerError).WithPayload(&models.Error{
					Data: &models.ErrorData{
						Message: err.Error(),
					},
				})
			}
		}
		return users.NewUserUpdateOK()
	}
}

func (c User) GetUserFunc() users.GetCurrentUserHandlerFunc {
	return func(p users.GetCurrentUserParams) middleware.Responder {
		user, err := c.client.User.Get(p.HTTPRequest.Context(), 1) // using 1 as id before auth implemented(jwt is not avaliable atm)
		id := int64(user.ID)
		if err != nil {
			return users.NewGetCurrentUserDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}

		return users.NewGetCurrentUserOK().WithPayload(&models.GetUserResponse{
			Data: &models.User{
				ID:                &id,
				Login:             &user.Login,
				Surname:           *user.Surname,
				Name:              user.Name,
				Patronymic:        *user.Patronymic,
				PassportSeries:    *user.PassportSeries,
				PassportNumber:    *user.PassportNumber,
				PassportAuthority: *user.PassportAuthority,
				PassportIssueDate: user.PassportIssueDate.String(),
				PhoneNumber:       *user.Phone,
				Email:             user.Email,
				Type:              string(user.Type),
				ActiveAreas:       user.ActiveAreas,
				OrgName:           *user.OrgName,
				Vk:                *user.Vk,
				Instagram:         *user.Instagram,
				Facebook:          *user.Facebook,
				Tiktok:            *user.Tiktok,
			},
		})
	}
}

func (c User) PatchUserFunc() users.PatchUserHandlerFunc {
	return func(p users.PatchUserParams, _ interface{}) middleware.Responder {
		return users.NewPatchUserNoContent()
	}
}

func (c User) AssignRoleToUserFunc() users.AssignRoleToUserHandlerFunc {
	return func(p users.AssignRoleToUserParams, access interface{}) middleware.Responder {
		_, err := authentication.IsAdmin(access)
		if err != nil {
			return users.NewAssignRoleToUserDefault(http.StatusInternalServerError).WithPayload(buildErrorPayload(err))
		}
		ctx := p.HTTPRequest.Context()
		userId := int(p.UserID)
		roleId := int(*p.Data.RoleID)
		foundUser, err := c.client.User.Get(ctx, userId)
		if err != nil {
			return users.NewAssignRoleToUserDefault(http.StatusNotFound).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		role, err := c.client.Role.Get(ctx, roleId)
		if err != nil {
			return users.NewAssignRoleToUserDefault(http.StatusNotFound).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		foundUser, err = c.client.User.UpdateOne(foundUser).SetRole(role).Save(ctx)
		if err != nil {
			return users.NewAssignRoleToUserDefault(http.StatusNotFound).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		userIdInt64 := int64(foundUser.ID)
		roleIdInt64 := int64(role.ID)
		return users.NewAssignRoleToUserOK().WithPayload(&models.GetUserResponse{
			Data: &models.User{
				CreateTime: nil,
				ID:         &userIdInt64,
				RoleID:     &roleIdInt64,
				Login:      &foundUser.Login,
			},
		})
	}
}
