package handlers

import (
	"net/http"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
)

// Keep in mind that these messages will be used by frontend team.
// So, if you change it, you should notify them.
const (
	allOk = "all ok"

	// Area

	errQueryTotalAreas = "failed to query total active areas"
	errQueryAreas      = "failed to query active areas"

	// Category

	errCreateCategory       = "cant create new category"
	errQueryTotalCategories = "cant get total amount of categories"
	errQueryCategories      = "cant get all categories"
	errGetCategory          = "failed to get category"
	errDeleteCategory       = "delete category failed"
	errUpdateCategory       = "cant update category"
	categoryDeleted         = "category deleted"

	// Email Confirmation

	errEmailConfirm = "failed to verify email confirmation token"
	emailConfirmed  = "you have successfully confirmed new email"

	// Equipment Periods

	errGetUnavailableEqStatus = "can't find unavailable equipment status dates"

	// Equipment Status Names

	errCreateEqStatus  = "can't create equipment status"
	errQueryEqStatuses = "can't get equipment statuses"
	errGetEqStatus     = "can't get equipment status"
	errDeleteEqStatus  = "can't delete equipment status"

	// Equipment Status

	errWrongEqStatus            = "wrong new equipment status, status should be only 'not available'"
	errGetEqStatusByID          = "can't find equipment status by provided id"
	errOrderAndUserByEqStatusID = "can't receive order and user data during checking equipment status"
	errUpdateEqStatus           = "can't update equipment status"

	// Equipment

	errCreateEquipment         = "error while creating equipment"
	errMapEquipment            = "error while mapping equipment"
	errGetEquipment            = "error while getting equipment"
	errEquipmentNotFound       = "equipment not found"
	errEquipmentArchive        = "error while archiving equipment"
	errEquipmentBlock          = "error while blocking equipment"
	errDeleteEquipment         = "error while deleting equipment"
	errQueryTotalEquipments    = "error while getting total of all equipments"
	errQueryEquipments         = "error while getting all equipments"
	errUpdateEquipment         = "error while updating equipment"
	errFindEquipment           = "error while finding equipment"
	errEquipmentBlockForbidden = "you don't have rights to block the equipment"
	equipmentDeleted           = "equipment deleted"

	// Order Status

	errQueryOrderHistory                 = "can't get order history"
	errQueryOrderHistoryForbidden        = "you don't have rights to see this order"
	errCreateOrderStatusForbidden        = "you don't have rights to add a new status"
	errOrderStatusEmpty                  = "order status is empty"
	errGetOrderStatus                    = "can't get order current status"
	errUpdateOrderStatus                 = "can't update status"
	errQueryTotalOrdersByStatus          = "can't get total count of orders by status"
	errQueryOrdersByStatus               = "can't get orders by status"
	errQueryTotalOrdersByPeriodAndStatus = "can't get total amount of orders by period and status"
	errQueryOrdersByPeriodAndStatus      = "can't get orders by period and status"
	errMapOrderStatus                    = "can't map order status name"
	errQueryStatusNames                  = "can't get all status names"

	// Order

	errOrderNotFound       = "no order with such id"
	errMapOrder            = "can't map order"
	errQueryOrders         = "can't get orders"
	errQueryTotalOrders    = "error while getting total of orders"
	errUpdateOrder         = "update order failed"
	errEquipmentIsNotFree  = "requested equipment is not free"
	errCheckEqStatusFailed = "error while checking if equipment is available for period"

	// Password Reset

	errLoginRequired       = "login is required"
	passwordResetSuccesful = "check your email for a reset link"

	// Pet Kind

	errCreatePetKind   = "error while creating pet kind"
	errGetPetKind      = "error while getting pet kind"
	errPetKindNotFound = "no pet kind found"
	errUpdatePetKind   = "error while updating pet kind"
	errDeletePetKind   = "error while deleting pet kind"
	petKindDeleted     = "pet kind deleted"

	// Pet Size

	errCreatePetSize        = "error while creating pet size"
	errPetSizeAlreadyExists = "error while creating pet size: the name already exist"
	errGetPetSize           = "error while getting pet size"
	errPetSizeNotFound      = "no pet size found"
	errUpdatePetSize        = "error while updating pet size"
	errDeletePetSize        = "error while deleting pet size"
	petSizeDeleted          = "pet size deleted"

	// Photo

	errCreatePhoto = "failed to save photo"
	errFileEmpty   = "File is empty"
	errWrongFormat = "Wrong file format. File should be jpg or jpeg"
	errGetPhoto    = "failed to get photo"
	errDeletePhoto = "failed to delete photo"
	photoDeleted   = "photo deleted"

	// Registration Confirm

	errRegistrationAlreadyConfirmed = "Registration is already confirmed."
	errRegistrationCannotFindUser   = "Can't find this user, registration confirmation link wasn't send"
	errRegistrationCannotSend       = "Can't send registration confirmation link. Please try again later"
	errFailedToConfirm              = "Failed to verify confirmation token. Please try again later"
	confirmationNotRequired         = "Confirmation link was not sent to email, sending parameter was set to false and not required"
	confirmationSent                = "Confirmation link was sent"
	registrationConfirmed           = "You have successfully confirmed registration"

	// Roles

	errQueryRoles = "can't get all roles"

	// Subcategory

	errCreateSubcategory   = "failed to create new subcategory"
	errMapSubcategory      = "failed to map new subcategory"
	errQuerySCatByCategory = "failed to list subcategories by category id"
	errGetSubcategory      = "failed to get subcategory"
	errDeleteSubcategory   = "failed to delete subcategory"
	errUpdateSubcategory   = "failed to update subcategory"
	subcategoryDeleted     = "subcategory deleted"

	// User

	errInvalidLoginOrPass   = "invalid login or password"
	errLoginInUse           = "login is already used"
	errCreateUser           = "failed to create user"
	errInvalidToken         = "token invalid"
	errTokenRefresh         = "error while refreshing token"
	errMapUser              = "map user error"
	errUserNotFound         = "can't find user by id"
	errUpdateUser           = "can't update user"
	errRoleRequired         = "role id is required"
	errSetUserRole          = "set user role error"
	errQueryTotalUsers      = "failed get user total amount"
	errQueryUsers           = "failed to get user list"
	errDeleteUser           = "can't delete user"
	errDeleteUserNotRO      = "user must be readonly for deletion"
	errUserPasswordChange   = "error while changing password"
	errWrongPassword        = "wrong password"
	errPasswordsAreSame     = "old and new passwords are the same"
	errPasswordPatchEmpty   = "password patch is empty"
	errUpdateROAccess       = "error while updating readonly access"
	errChangeEmail          = "error while changing email"
	errEmailPatchEmpty      = "email patch is empty"
	errNewEmailConfirmation = "can't send link for confirmation new email"
	logoutSuccessful        = "successfully logged out"
	roleAssigned            = "role assigned"
)

func buildErrorPayload(code int32, msg string, details string) *models.SwaggerError {
	return &models.SwaggerError{
		Code:    &code,
		Message: &msg,
		Details: details, // optional field for raw err messages
	}
}

func buildInternalErrorPayload(msg string, details string) *models.SwaggerError {
	return buildErrorPayload(http.StatusInternalServerError, msg, details)
}

func buildExFailedErrorPayload(msg string, details string) *models.SwaggerError {
	return buildErrorPayload(http.StatusExpectationFailed, msg, details)
}

func buildConflictErrorPayload(msg string, details string) *models.SwaggerError {
	return buildErrorPayload(http.StatusConflict, msg, details)
}

func buildNotFoundErrorPayload(msg string, details string) *models.SwaggerError {
	return buildErrorPayload(http.StatusNotFound, msg, details)
}

func buildForbiddenErrorPayload(msg string, details string) *models.SwaggerError {
	return buildErrorPayload(http.StatusForbidden, msg, details)
}

func buildBadRequestErrorPayload(msg string, details string) *models.SwaggerError {
	return buildErrorPayload(http.StatusBadRequest, msg, details)
}
