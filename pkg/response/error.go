package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Errors struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

func (e *Errors) Error() string {
	return e.Err.Error()
}

func NewError(code int, err string) Errors {
	return Errors{
		Code: code,
		Err:  errors.New(err),
	}
}

var (
	// Database Errors
	ErrConnectDatabase = NewError(fiber.StatusInternalServerError, "Unable to connect to the database. Please try again later.")
	ErrMigrateDatabase = NewError(fiber.StatusInternalServerError, "There was an issue with database migration. Please contact support.")

	// User Errors
	ErrUserNotFound           = NewError(fiber.StatusNotFound, "User not found. Please check your details and try again.")
	ErrUserAlreadyExists      = NewError(fiber.StatusConflict, "This email is already registered. Please try logging in or use a different email.")
	ErrHashPassword           = NewError(fiber.StatusInternalServerError, "Something went wrong while processing your password. Please try again.")
	ErrGenerateToken          = NewError(fiber.StatusInternalServerError, "There was an error generating your authentication token. Please try again.")
	ErrInvalidEmail           = NewError(fiber.StatusBadRequest, "Please provide a valid email address.")
	ErrInvalidPassword        = NewError(fiber.StatusBadRequest, "The password you entered is incorrect. Please try again.")
	ErrJWTToken               = NewError(fiber.StatusInternalServerError, "Something went wrong with the token. Please try again.")
	ErrFailedSendNotification = NewError(fiber.StatusInternalServerError, "There was an error sending the notification. Please try again later.")
	ErrUnauthorized           = NewError(fiber.StatusUnauthorized, "Invalid login credentials. Please check your username and password.")
	ErrSetHTML                = NewError(fiber.StatusInternalServerError, "There was an error setting the HTML template. Please try again.")
	ErrExecuteHTML            = NewError(fiber.StatusInternalServerError, "We encountered an issue while processing the page. Please try again.")
	ErrCredentialMismatch     = NewError(fiber.StatusUnauthorized, "The credentials provided do not match our records.")
	ErrForbiddenRole          = NewError(fiber.StatusForbidden, "You do not have permission to access this resource.")
	ErrFailedDeleteImage      = NewError(fiber.StatusInternalServerError, "Failed to delete the image. Please try again.")

	// Google Errors
	ErrStateNoMatch                = NewError(fiber.StatusBadRequest, "The state does not match. Please try again.")
	ErrFailedFetchUserInfo         = NewError(fiber.StatusInternalServerError, "Failed to fetch your Google user data. Please try again later.")
	ErrReadResponseBody            = NewError(fiber.StatusInternalServerError, "There was an issue reading the response. Please try again.")
	ErrUnmarshal                   = NewError(fiber.StatusBadRequest, "There was a problem with the data we received. Please try again.")
	ErrRegisterInternalServerError = NewError(fiber.StatusInternalServerError, "An error occurred while registering. Please try again later.")
	ErrInvalidCode                 = NewError(fiber.StatusBadRequest, "The code provided is invalid. Please check and try again.")

	// Chat Errors
	ErrChatFailed     = NewError(fiber.StatusInternalServerError, "Failed to create chat. Please try again.")
	ErrSaveChatFailed = NewError(fiber.StatusInternalServerError, "Failed to save chat. Please try again.")
)
