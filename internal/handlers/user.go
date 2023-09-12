package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/user"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations/users"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/repositories"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/utils"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/pkg/domain"
)

func SetUserHandler(logger *zap.Logger, api *operations.BeAPI,
	tokenManager domain.TokenManager,
	regConfirmService domain.RegistrationConfirmService, changeEmailService domain.ChangeEmailService) {
	userRepo := repositories.NewUserRepository()
	userHandler := NewUser(logger)

	api.UsersLoginHandler = userHandler.LoginUserFunc(tokenManager)
	api.UsersRefreshHandler = userHandler.Refresh(tokenManager)
	api.UsersPostUserHandler = userHandler.PostUserFunc(userRepo, regConfirmService)
	api.UsersGetCurrentUserHandler = userHandler.GetUserFunc(userRepo)
	api.UsersPatchUserHandler = userHandler.PatchUserFunc(userRepo)
	api.UsersGetUserHandler = userHandler.GetUserById(userRepo)
	api.UsersGetAllUsersHandler = userHandler.GetUsersList(userRepo)
	api.UsersAssignRoleToUserHandler = userHandler.AssignRoleToUserFunc(userRepo)
	api.UsersChangePasswordHandler = userHandler.ChangePassword(userRepo)
	api.UsersLogoutHandler = userHandler.LogoutUserFunc(tokenManager)
	api.UsersDeleteCurrentUserHandler = userHandler.DeleteCurrentUser(userRepo)
	api.UsersDeleteUserHandler = userHandler.DeleteUser(userRepo)
	api.UsersUpdateReadonlyAccessHandler = userHandler.UpdateReadonlyAccess(userRepo)
	api.UsersChangeEmailHandler = userHandler.ChangeEmail(userRepo, changeEmailService)
}

type User struct {
	logger *zap.Logger
}

func NewUser(logger *zap.Logger) *User {
	return &User{
		logger: logger,
	}
}

func (c User) LoginUserFunc(service domain.TokenManager) users.LoginHandlerFunc {
	return func(p users.LoginParams) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		login := *p.Login.Login
		password := *p.Login.Password
		accessToken, refreshToken, isInternalErr, err := service.GenerateTokens(ctx, login, password)
		if err != nil {
			if isInternalErr {
				return users.NewLoginDefault(http.StatusInternalServerError)
			}
			return users.NewLoginUnauthorized().WithPayload("Invalid login or password")
		}

		return users.NewLoginOK().WithPayload(&models.TokenPair{
			AccessToken:  &accessToken,
			RefreshToken: &refreshToken,
		})
	}
}

func (c User) LogoutUserFunc(tokenManager domain.TokenManager) users.LogoutHandlerFunc {
	return func(p users.LogoutParams) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		refreshToken := *p.RefreshToken.RefreshToken
		err := tokenManager.DeleteTokenPair(ctx, refreshToken)
		if err != nil && ent.IsNotFound(err) {
			return users.NewLogoutNotFound()
		}
		if err != nil {
			return users.NewLogoutDefault(http.StatusInternalServerError)
		}
		return users.NewLogoutOK().WithPayload("Successfully logged out")
	}
}

func (c User) PostUserFunc(repository domain.UserRepository, regConfirmService domain.RegistrationConfirmService) users.PostUserHandlerFunc {
	return func(p users.PostUserParams) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		createdUser, err := repository.CreateUser(ctx, p.Data)
		if err != nil {
			if ent.IsConstraintError(err) {
				return users.NewPostUserDefault(http.StatusExpectationFailed).
					WithPayload(buildExFailedErrorPayload("login is already used"))
			}
			return users.NewPostUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(err.Error()))
		}

		id := int64(createdUser.ID)

		err = regConfirmService.SendConfirmationLink(ctx, createdUser.Login)
		if err != nil {
			c.logger.Error("error sending registration confirmation link", zap.Error(err))
		}

		return users.NewPostUserCreated().WithPayload(&models.CreateUserResponse{
			Data: &models.CreateUserResponseData{
				ID:    &id,
				Login: &createdUser.Login,
			},
		})
	}
}

func (c User) Refresh(manager domain.TokenManager) users.RefreshHandlerFunc {
	return func(p users.RefreshParams) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		refreshToken := *p.RefreshToken.RefreshToken
		newAccess, NewRefresh, isValid, err := manager.RefreshToken(ctx, refreshToken)
		if isValid {
			c.logger.Info("token invalid", zap.String("token", refreshToken))
			return users.NewRefreshDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload("token invalid"))
		}
		if err != nil {
			c.logger.Error("Error while refreshing token", zap.Error(err))
			return users.NewRefreshDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while refreshing token"))
		}
		return users.NewRefreshOK().WithPayload(&models.TokenPair{
			AccessToken:  &newAccess,
			RefreshToken: &NewRefresh,
		})
	}
}

func (c User) GetUserFunc(repository domain.UserRepository) users.GetCurrentUserHandlerFunc {
	return func(p users.GetCurrentUserParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)

		user, err := repository.GetUserByID(ctx, userID)
		if err != nil {
			c.logger.Error("get user by id error", zap.Error(err))
			return users.NewGetCurrentUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("cant find user by id"))
		}

		result, err := mapUserInfo(user)
		if err != nil {
			c.logger.Error("map user error", zap.Error(err))
			return users.NewGetCurrentUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("map user error"))
		}

		return users.NewGetCurrentUserOK().WithPayload(result)
	}
}

func (c User) PatchUserFunc(repository domain.UserRepository) users.PatchUserHandlerFunc {
	return func(p users.PatchUserParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)

		err := repository.UpdateUserByID(ctx, userID, p.UserPatch)
		if err != nil {
			c.logger.Error("get user by id error", zap.Error(err))
			return users.NewPatchUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("cant update user"))
		}
		return users.NewPatchUserNoContent()
	}
}

func (c User) AssignRoleToUserFunc(repository domain.UserRepository) users.AssignRoleToUserHandlerFunc {
	return func(p users.AssignRoleToUserParams, principal *models.Principal) middleware.Responder {

		ctx := p.HTTPRequest.Context()
		userId := int(p.UserID)
		if p.Data.RoleID == nil {
			return users.NewAssignRoleToUserDefault(http.StatusBadRequest).
				WithPayload(buildInternalErrorPayload("role id is required"))
		}
		roleId := int(*p.Data.RoleID)

		err := repository.SetUserRole(ctx, userId, roleId)
		if err != nil {
			c.logger.Error("set user role error", zap.Error(err))
			return users.NewAssignRoleToUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload(err.Error()))
		}

		return users.NewAssignRoleToUserOK().WithPayload("role assigned")
	}
}

func (c User) GetUserById(repository domain.UserRepository) users.GetUserHandlerFunc {
	return func(p users.GetUserParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		id := int(p.UserID)
		foundUser, err := repository.GetUserByID(ctx, id)
		if err != nil {
			return users.NewGetUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("cant find user by id"))
		}

		userToResponse, err := mapUserInfo(foundUser)
		if err != nil {
			c.logger.Error("map user error", zap.Error(err))
			return users.NewGetUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("map user error"))
		}

		return users.NewGetUserOK().WithPayload(userToResponse)
	}
}

func (c User) GetUsersList(repository domain.UserRepository) users.GetAllUsersHandlerFunc {
	return func(p users.GetAllUsersParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		limit := utils.GetValueByPointerOrDefaultValue(p.Limit, math.MaxInt)
		offset := utils.GetValueByPointerOrDefaultValue(p.Offset, 0)
		orderBy := utils.GetValueByPointerOrDefaultValue(p.OrderBy, utils.AscOrder)
		orderColumn := utils.GetValueByPointerOrDefaultValue(p.OrderColumn, user.FieldID)
		total, err := repository.UsersListTotal(ctx)
		if err != nil {
			c.logger.Error("failed get user total amount", zap.Error(err))
			return users.NewGetAllUsersDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("failed to get user total amount"))
		}
		var allUsers []*ent.User
		if total > 0 {
			allUsers, err = repository.UserList(ctx, int(limit), int(offset), orderBy, orderColumn)
			if err != nil {
				c.logger.Error("failed get user list", zap.Error(err))
				return users.NewGetAllUsersDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload("failed to get user list"))
			}
		}

		usersToResponse := make([]*models.UserInfo, len(allUsers))
		for i, element := range allUsers {
			userToResponse, errMap := mapUserInfo(element)
			if errMap != nil {
				c.logger.Error("map user error", zap.Error(errMap))
				return users.NewGetAllUsersDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload("map user error"))
			}
			usersToResponse[i] = userToResponse
		}
		totalUsers := int64(total)
		listUsers := &models.GetListUsers{
			Items: usersToResponse,
			Total: &totalUsers,
		}

		return users.NewGetAllUsersOK().WithPayload(listUsers)
	}
}

func (c User) DeleteCurrentUser(repository domain.UserRepository) users.DeleteCurrentUserHandlerFunc {
	return func(p users.DeleteCurrentUserParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)

		err := repository.Delete(ctx, userID)
		if err != nil {
			c.logger.Error("error during deleting user", zap.Error(err))
			return users.NewDeleteCurrentUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("can't delete user"))
		}

		return users.NewDeleteCurrentUserOK()
	}
}

func (c User) DeleteUser(repo domain.UserRepository) users.DeleteUserHandlerFunc {
	return func(p users.DeleteUserParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(p.UserID)
		deletedByUserID := int(principal.ID)

		user, err := repo.GetUserByID(ctx, userID)
		if err != nil {
			c.logger.Error(fmt.Sprintf("retrieving user by ID %d", userID), zap.Error(err))
			if ent.IsNotFound(err) {
				return users.NewDeleteUserNotFound().
					WithPayload(buildNotFoundErrorPayload("User not found"))
			}
			return users.NewDeleteUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Unexpected error"))
		}

		if !user.IsReadonly {
			c.logger.Error("User must be readonly for deletion", zap.Int("userID", userID))
			return users.NewDeleteUserDefault(http.StatusForbidden).
				WithPayload(buildForbiddenErrorPayload("User must be readonly for deletion"))
		}

		if err := repo.Delete(ctx, userID); err != nil {
			c.logger.Error("deleting user", zap.Error(err))
			if ent.IsNotFound(err) {
				return users.NewDeleteUserNotFound().
					WithPayload(buildNotFoundErrorPayload("User not found"))
			}
			return users.NewDeleteUserDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Unexpected error"))
		}

		c.logger.Info("User deleted successfully", zap.Int("userID", userID), zap.Int("deletedByUserID", deletedByUserID))
		return users.NewDeleteUserNoContent()
	}
}

func (c User) ChangePassword(repo domain.UserRepository) users.ChangePasswordHandlerFunc {
	return func(p users.ChangePasswordParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)

		if p.PasswordPatch == nil {
			c.logger.Error("password patch is nil", zap.Any("principal", principal))
			return users.NewChangePasswordDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload("Password patch is nil"))
		}
		//TODO: add validation for password or ask frontend to do it
		if p.PasswordPatch.OldPassword == p.PasswordPatch.NewPassword {
			c.logger.Error("old and new passwords are the same", zap.Any("principal", principal))
			return users.NewChangePasswordDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload("Old and new passwords are the same"))
		}
		requestedUser, err := repo.GetUserByID(ctx, userID)
		if err != nil {
			c.logger.Error("getting user failed", zap.Error(err))
			return users.NewChangePasswordDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Can't get user by id"))
		}
		expectedPasswordHash := requestedUser.Password
		if err = bcrypt.CompareHashAndPassword([]byte(expectedPasswordHash), []byte(p.PasswordPatch.OldPassword)); err != nil {
			c.logger.Error("wrong password", zap.Error(err))
			return users.NewChangePasswordDefault(http.StatusForbidden).
				WithPayload(buildForbiddenErrorPayload("Wrong password"))
		}
		if err = repo.ChangePasswordByLogin(ctx, requestedUser.Login, p.PasswordPatch.NewPassword); err != nil {
			c.logger.Error("error while changing password", zap.Error(err))
			return users.NewChangePasswordDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while changing password"))
		}
		return users.NewChangePasswordNoContent()
	}
}

func (c User) UpdateReadonlyAccess(repo domain.UserRepository) users.UpdateReadonlyAccessHandlerFunc {
	return func(p users.UpdateReadonlyAccessParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		currentUserID := int(principal.ID)
		userID := int(p.UserID)
		isReadonly := p.Body.IsReadonly

		if err := repo.SetIsReadonly(ctx, userID, isReadonly); err != nil {
			c.logger.Error("error while updating readonly access", zap.Error(err))
			if ent.IsNotFound(err) {
				return users.NewUpdateReadonlyAccessNotFound().
					WithPayload(buildNotFoundErrorPayload("User not found"))
			}
			return users.NewUpdateReadonlyAccessDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Unexpected error"))
		}

		if isReadonly {
			c.logger.Info(fmt.Sprintf("User %d has been granted read-only access by user %d", userID, currentUserID))
		} else {
			c.logger.Info(fmt.Sprintf("Read-only access for user %d has been revoked by user %d", userID, currentUserID))
		}

		return users.NewUpdateReadonlyAccessNoContent()
	}
}

func (c User) ChangeEmail(repo domain.UserRepository,
	changeEmailService domain.ChangeEmailService) users.ChangeEmailHandlerFunc {
	return func(p users.ChangeEmailParams, principal *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		userID := int(principal.ID)

		if p.EmailPatch == nil {
			c.logger.Error("email patch is nil", zap.Any("principal", principal))
			return users.NewChangeEmailDefault(http.StatusBadRequest).
				WithPayload(buildBadRequestErrorPayload("Email patch is nil"))
		}

		requestedUser, err := repo.GetUserByID(ctx, userID)
		if err != nil {
			c.logger.Error("getting user for changing email failed", zap.Error(err))
			return users.NewChangeEmailDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Can't get user by id"))
		}

		err = changeEmailService.SendEmailConfirmationLink(ctx, requestedUser.Login, p.EmailPatch.NewEmail)
		if err != nil {
			c.logger.Error("error while sending link for confirmation new email", zap.Error(err))
			return users.NewChangeEmailDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Can't send link for confirmation new email"))
		}

		return users.NewChangeEmailNoContent()
	}
}

func mapUserInfo(user *ent.User) (*models.UserInfo, error) {
	userID := int64(user.ID)
	passportDate := user.PassportIssueDate.String()
	if user.Edges.Role == nil {
		return nil, errors.New("role is nil")
	}
	userRole := user.Edges.Role
	userRoleInfo := models.UserInfoRole{
		ID:   int64(userRole.ID),
		Name: userRole.Name,
	}
	typeString := user.Type.String()
	result := &models.UserInfo{
		Email:                   &user.Email,
		ID:                      &userID,
		IsReadonly:              &user.IsReadonly,
		Login:                   &user.Login,
		Name:                    &user.Name,
		OrgName:                 user.OrgName,
		PassportAuthority:       user.PassportAuthority,
		PassportIssueDate:       &passportDate,
		PassportNumber:          user.PassportNumber,
		PassportSeries:          user.PassportSeries,
		Patronymic:              user.Patronymic,
		PhoneNumber:             user.Phone,
		Role:                    &userRoleInfo,
		Surname:                 user.Surname,
		Type:                    &typeString,
		IsRegistrationConfirmed: &user.IsRegistrationConfirmed,
	}
	return result, nil
}
