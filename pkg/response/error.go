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
	ErrConnectDatabase = NewError(fiber.StatusInternalServerError, "We are experiencing issues connecting to our database. Please try again later.")
	ErrMigrateDatabase = NewError(fiber.StatusInternalServerError, "We encountered a problem with our database setup. Please contact support for assistance.")

	// User Errors
	ErrUserNotFound           = NewError(fiber.StatusNotFound, "We couldn't find a user with the provided details. Please check and try again.")
	ErrUserAlreadyExists      = NewError(fiber.StatusConflict, "This email is already registered. Please log in or use a different email.")
	ErrHashPassword           = NewError(fiber.StatusInternalServerError, "We encountered an issue processing your password. Please try again.")
	ErrGenerateToken          = NewError(fiber.StatusInternalServerError, "We had trouble generating your authentication token. Please try again.")
	ErrInvalidEmail           = NewError(fiber.StatusBadRequest, "The email address you provided is invalid. Please enter a valid email.")
	ErrInvalidPassword        = NewError(fiber.StatusBadRequest, "The password you entered is incorrect. Please try again.")
	ErrJWTToken               = NewError(fiber.StatusInternalServerError, "We encountered an issue with your token. Please try again.")
	ErrFailedSendNotification = NewError(fiber.StatusInternalServerError, "We couldn't send the notification. Please try again later.")
	ErrUnauthorized           = NewError(fiber.StatusUnauthorized, "Your login credentials are incorrect. Please check and try again.")
	ErrSetHTML                = NewError(fiber.StatusInternalServerError, "We encountered an issue setting up the page. Please try again.")
	ErrExecuteHTML            = NewError(fiber.StatusInternalServerError, "We had trouble processing the page. Please try again.")
	ErrCredentialMismatch     = NewError(fiber.StatusUnauthorized, "The credentials you provided don't match our records. Please try again.")
	ErrForbiddenRole          = NewError(fiber.StatusForbidden, "You don't have permission to access this resource.")
	ErrFailedDeleteImage      = NewError(fiber.StatusInternalServerError, "We couldn't delete the image. Please try again.")

	// Chat Errors
	ErrChatFailed     = NewError(fiber.StatusInternalServerError, "We couldn't create the chat. Please try again.")
	ErrSaveChatFailed = NewError(fiber.StatusInternalServerError, "We couldn't save the chat. Please try again.")

	// Payment Errors
	ErrSavePayment         = NewError(fiber.StatusInternalServerError, "We encountered an issue processing your payment. Please try again.")
	ErrUpdateStatus        = NewError(fiber.StatusInternalServerError, "We couldn't update the status. Please try again.")
	ErrFetchStatus         = NewError(fiber.StatusInternalServerError, "We couldn't retrieve the status. Please try again.")
	ErrFailedFindUser      = NewError(fiber.StatusInternalServerError, "We couldn't find the user. Please try again.")
	ErrFailedCreatePayment = NewError(fiber.StatusInternalServerError, "We couldn't process your payment. Please try again.")
	ErrMissingStatus       = NewError(fiber.StatusBadRequest, "The payment status is missing. Please provide a valid status.")
)
